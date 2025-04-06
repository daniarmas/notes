package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/daniarmas/clogg"
	"github.com/daniarmas/notes/internal/cache"
	"github.com/daniarmas/notes/internal/config"
	"github.com/daniarmas/notes/internal/data"
	"github.com/daniarmas/notes/internal/database"
	"github.com/daniarmas/notes/internal/domain"
	"github.com/daniarmas/notes/internal/httpserver"
	"github.com/daniarmas/notes/internal/k8sc"
	"github.com/daniarmas/notes/internal/oss"
	"github.com/daniarmas/notes/internal/service"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/cobra"
)

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Set up clogg
	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := clogg.GetLogger(clogg.LoggerConfig{
		BufferSize: 100,
		Handler:    handler,
	})
	defer logger.Shutdown()

	// Kubernetes client
	k8sClient, ck8sError := k8sc.NewClient()
	if ck8sError != nil {
		clogg.Error(ctx, "error creating k8s client", clogg.String("error", ck8sError.Error()))
	}

	// Config
	cfg := config.LoadServerConfig()

	// Set if running in k8s
	if k8sClient != nil {
		clogg.Info(ctx, "app is running in k8s")
		cfg.InK8s = true
	}

	// Database connection
	db := database.Open(ctx, cfg, true)
	defer database.Close(ctx, db, true)

	// Database queries
	dbQueries := database.New(db)

	// Cache connection
	rdb := cache.OpenRedis(ctx, cfg)
	defer cache.CloseRedis(ctx, rdb)

	// Object storage service
	oss := oss.NewDigitalOceanWithMinio(cfg)
	// Healthcheck
	if err := oss.HealthCheck(); err != nil {
		clogg.Error(ctx, "oss service healthcheck failed", clogg.String("error", err.Error()))
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
	authenticationService := service.NewAuthenticationService(jwtDatasource, hashDatasource, userRepository, accessTokenRepository, refreshTokenRepository, db)
	noteService := service.NewNoteService(noteRepository, oss, fileRepository, *cfg, k8sClient, db)

	// Http server
	resSrv := httpserver.NewRestServer(authenticationService, noteService, jwtDatasource, *cfg)
	graphSrv := httpserver.NewGraphQLServer(authenticationService, noteService, *cfg, jwtDatasource)

	// Start the Rest server
	go func() {
		msg := fmt.Sprintf("Http server listening on %s\n", resSrv.HttpServer.Addr)
		clogg.Info(ctx, msg)
		if err := resSrv.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			clogg.Error(ctx, "error listening and serving rest server", clogg.String("error", err.Error()))
		}
	}()

	// Start the GraphQL server
	go func() {
		msg := fmt.Sprintf("connect to http://localhost:%s/ for GraphQL playground", graphSrv.GraphQLServer.Addr)
		clogg.Info(ctx, msg)
		if err := graphSrv.GraphQLServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			clogg.Error(ctx, "error listening and serving graphql server", clogg.String("error", err.Error()))
		}
	}()

	// Graceful shutdown
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := resSrv.HttpServer.Shutdown(shutdownCtx); err != nil {
			clogg.Error(ctx, "error shutting down http server", clogg.String("error", err.Error()))
		}
		if err := graphSrv.GraphQLServer.Shutdown(shutdownCtx); err != nil {
			clogg.Error(ctx, "error shutting down graphql server", clogg.String("error", err.Error()))
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
