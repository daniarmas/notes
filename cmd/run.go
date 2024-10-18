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
	"github.com/daniarmas/notes/internal/service"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/cobra"
)

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Custom logger
	clog.NewClog()

	// Config
	cfg := config.LoadConfig()

	// Database connection
	db := database.Open(cfg, true)
	defer database.Close(db, true)

	// Database queries
	dbQueries := database.New(db)

	// Cache connection
	rdb := cache.OpenRedis(cfg)
	defer cache.CloseRedis(rdb)

	// Datasources
	hashDatasource := data.NewBcryptHashDatasource()
	jwtDatasource := domain.NewJWTDatasource(cfg)
	userCacheDs := data.NewUserCacheDs(rdb)
	userDatabaseDs := data.NewUserDatabaseDs(dbQueries)
	accessTokenCacheDs := data.NewAccessTokenTokenCacheDs(rdb)
	accessTokenDatabaseDs := data.NewAccessTokenDatabaseDs(dbQueries)
	refreshTokenCacheDs := data.NewRefreshTokenCacheDs(rdb)
	refreshTokenDatabaseDs := data.NewRefreshTokenDatabaseDs(dbQueries)

	// Repositories
	userRepository := domain.NewUserRepository(&userCacheDs, &userDatabaseDs)
	accessTokenRepository := domain.NewAccessTokenRepository(accessTokenCacheDs, accessTokenDatabaseDs)
	refreshTokenRepository := domain.NewRefreshTokenRepository(&refreshTokenCacheDs, &refreshTokenDatabaseDs)

	// Services
	authenticationService := service.NewAuthenticationService(jwtDatasource, hashDatasource, userRepository, accessTokenRepository, refreshTokenRepository)

	// Http server
	srv := httpserver.NewServer(authenticationService, jwtDatasource)

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
