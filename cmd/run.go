package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/daniarmas/notes/internal/cache"
	"github.com/daniarmas/notes/internal/clog"
	"github.com/daniarmas/notes/internal/config"
	"github.com/daniarmas/notes/internal/data"
	"github.com/daniarmas/notes/internal/database"
	"github.com/daniarmas/notes/internal/domain"
	"github.com/daniarmas/notes/internal/httpserver"
	"github.com/daniarmas/notes/internal/oss"
	"github.com/daniarmas/notes/internal/service"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Custom logger
	clog.NewClog()

	// Kubernetes client
	// creates the in-cluster config
	var k8sClient *kubernetes.Clientset
	k8sCfg, k8sCfgErr := rest.InClusterConfig()
	if k8sCfgErr != nil {
		clog.Error(ctx, "error creating in-cluster config", k8sCfgErr)
	}
	clog.Info(ctx, "Kubernetes in-cluster config loaded", nil)

	// K8s clientset
	if k8sCfgErr != nil {
		var err error
		k8sClient, err = kubernetes.NewForConfig(k8sCfg)
		if err != nil {
			clog.Error(ctx, "error creating kubernetes clientset", err)
		}
	}

	// Config
	cfg := config.LoadServerConfig()

	// Set if running in k8s
	if k8sCfg != nil {
		clog.Info(ctx, "app is running in k8s", nil)
		cfg.InK8s = true
	}

	// Database connection
	db := database.Open(cfg, true)
	defer database.Close(db, true)

	// Database queries
	dbQueries := database.New(db)

	// Cache connection
	rdb := cache.OpenRedis(cfg)
	defer cache.CloseRedis(rdb)

	// Object storage service
	oss := oss.NewDigitalOceanWithMinio(cfg)
	// Healthcheck
	if err := oss.HealthCheck(); err != nil {
		clog.Error(ctx, "error checking object storage service health", err)
	}

	// Datasources
	hashDatasource := data.NewBcryptHashDatasource()
	jwtDatasource := domain.NewJWTDatasource(cfg)
	userCacheDs := data.NewUserCacheDs(rdb)
	userDatabaseDs := data.NewUserDatabaseDs(dbQueries)
	accessTokenCacheDs := data.NewAccessTokenTokenCacheDs(rdb)
	accessTokenDatabaseDs := data.NewAccessTokenDatabaseDs(dbQueries)
	refreshTokenCacheDs := data.NewRefreshTokenCacheDs(rdb)
	refreshTokenDatabaseDs := data.NewRefreshTokenDatabaseDs(dbQueries)
	noteCacheDs := data.NewNoteCacheDs(rdb)
	noteDatabaseDs := data.NewNoteDatabaseDs(dbQueries)
	fileDatabaseDs := data.NewFileDatabaseDs(dbQueries)

	// Repositories
	userRepository := domain.NewUserRepository(&userCacheDs, &userDatabaseDs)
	accessTokenRepository := domain.NewAccessTokenRepository(accessTokenCacheDs, accessTokenDatabaseDs)
	refreshTokenRepository := domain.NewRefreshTokenRepository(&refreshTokenCacheDs, &refreshTokenDatabaseDs)
	noteRepository := domain.NewNoteRepository(&noteCacheDs, &noteDatabaseDs)
	fileRepository := domain.NewFileRepository(fileDatabaseDs, oss, cfg)

	// Services
	authenticationService := service.NewAuthenticationService(jwtDatasource, hashDatasource, userRepository, accessTokenRepository, refreshTokenRepository)
	noteService := service.NewNoteService(noteRepository, oss, fileRepository, *cfg, k8sClient)

	// Http server
	srv := httpserver.NewServer(authenticationService, noteService, jwtDatasource)

	go func() {
		msg := fmt.Sprintf("Http server listening on %s\n", srv.HttpServer.Addr)
		clog.Info(ctx, msg, nil)
		if err := srv.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			clog.Error(ctx, "error listening and serving", err)
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := srv.HttpServer.Shutdown(shutdownCtx); err != nil {
			clog.Error(ctx, "error shutting down http server", err)
		}
	}()
	wg.Wait()
	return nil
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		if err := run(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
