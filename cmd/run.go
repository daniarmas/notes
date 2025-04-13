package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/daniarmas/clogg"
	httpw "github.com/daniarmas/http"
	cmiddleware "github.com/daniarmas/http/middleware"
	"github.com/daniarmas/notes/internal/cache"
	"github.com/daniarmas/notes/internal/config"
	"github.com/daniarmas/notes/internal/data"
	"github.com/daniarmas/notes/internal/database"
	"github.com/daniarmas/notes/internal/domain"
	"github.com/daniarmas/notes/internal/httpserver"
	"github.com/daniarmas/notes/internal/httpserver/handler"
	"github.com/daniarmas/notes/internal/httpserver/middleware"
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
	jsonhandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := clogg.GetLogger(clogg.LoggerConfig{
		BufferSize: 100,
		Handler:    jsonhandler,
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
	db, err := database.Open(ctx, cfg)
	if err != nil {
		clogg.Error(ctx, "error opening database", clogg.String("error", err.Error()))
		os.Exit(1)
	}
	defer database.Close(ctx, db)

	// Database queries
	dbQueries := database.New(db)

	// Cache connection
	rdb, err := cache.OpenRedis(ctx, cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword, cfg.RedisDb)
	if err != nil {
		clogg.Error(ctx, "error connecting to redis", clogg.String("error", err.Error()))
		os.Exit(1)
	}
	clogg.Info(ctx, "connected to redis", clogg.String("host", cfg.RedisHost), clogg.String("port", cfg.RedisPort))

	// Close redis connection
	defer func() {
		err := rdb.Close()
		if err != nil {
			clogg.Error(ctx, "error closing redis connection", clogg.String("error", err.Error()))
			os.Exit(1)
		}
		clogg.Info(ctx, "redis connection closed")
	}()

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

	// Httpw server
	routes := []httpw.HandleFunc{
		// OpenAPI specification
		{Pattern: "GET /swagger.json", Handler: handler.OpenApiHanlder},
		// Authentication
		{Pattern: "GET /me", Handler: middleware.LoggedOnly(handler.Me(authenticationService)).(http.HandlerFunc)},
		{Pattern: "POST /sign-in", Handler: handler.SignIn(authenticationService)},
		{Pattern: "POST /sign-out", Handler: middleware.LoggedOnly(handler.SignOut(authenticationService)).(http.HandlerFunc)},
		// Note
		{Pattern: "GET /note/trash", Handler: middleware.LoggedOnly(handler.ListTrashNotesByUser(noteService)).(http.HandlerFunc)},
		{Pattern: "GET /note", Handler: middleware.LoggedOnly(handler.ListNotesByUser(noteService)).(http.HandlerFunc)},
		{Pattern: "POST /note", Handler: middleware.LoggedOnly(handler.CreateNote(noteService)).(http.HandlerFunc)},
		{Pattern: "DELETE /note/{id}/hard", Handler: middleware.LoggedOnly(handler.HardDeleteNote(noteService)).(http.HandlerFunc)},
		{Pattern: "DELETE /note/{id}", Handler: middleware.LoggedOnly(handler.SoftDeleteNote(noteService)).(http.HandlerFunc)},
		{Pattern: "PATCH /note/{id}", Handler: middleware.LoggedOnly(handler.UpdateNote(noteService)).(http.HandlerFunc)},
		{Pattern: "PATCH /note/{id}/restore", Handler: middleware.LoggedOnly(handler.RestoreNote(noteService)).(http.HandlerFunc)},
		{Pattern: "POST /note/presigned-urls", Handler: middleware.LoggedOnly(handler.GetPresignedUrls(noteService)).(http.HandlerFunc)},
	}

	httpwServer := httpw.New(httpw.Options{
		Addr:         net.JoinHostPort("0.0.0.0", cfg.RestServerPort),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
		Middlewares: []cmiddleware.Middleware{
			cmiddleware.LoggingMiddleware,
			middleware.SetUserInContext(jwtDatasource),
			cmiddleware.AllowCors(cmiddleware.CorsOptions{
				AllowedOrigin:  fmt.Sprintf("http://localhost:%s", cfg.RestServerPort),
				AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders: []string{"Content-Type", "Authorization"},
			}),
			cmiddleware.RecoverMiddleware,
		},
	}, routes...)

	// Http server
	graphSrv := httpserver.NewGraphQLServer(authenticationService, noteService, *cfg, jwtDatasource)

	var wg sync.WaitGroup
	wg.Add(3)

	// Start the rest server in a separate goroutine.
	go func() {
		defer wg.Done()
		if err := httpwServer.Run(ctx); err != nil {
			clogg.Error(ctx, "error listening and serving rest server", clogg.String("error", err.Error()))
		}
	}()

	// Start the GraphQL server
	go func() {
		defer wg.Done()
		msg := fmt.Sprintf("connect to http://localhost:%s/ for GraphQL playground", graphSrv.GraphQLServer.Addr)
		clogg.Info(ctx, msg)
		if err := graphSrv.GraphQLServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			clogg.Error(ctx, "error listening and serving graphql server", clogg.String("error", err.Error()))
		}
	}()

	// Graceful shutdown
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		// Shutdown the rest server
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
