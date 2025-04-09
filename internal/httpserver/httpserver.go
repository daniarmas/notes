package httpserver

import (
	"net"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/daniarmas/notes/internal/config"
	"github.com/daniarmas/notes/internal/domain"
	"github.com/daniarmas/notes/internal/graph"
	"github.com/daniarmas/notes/internal/httpserver/middleware"
	"github.com/daniarmas/notes/internal/service"
)

type HandleFunc struct {
	Pattern string
	Handler http.HandlerFunc
}

type Server struct {
	Mux           *http.ServeMux
	GraphQLServer *http.Server
}

// NewGraphQLServer creates and configures a new GraphQL server with the specified address.
func NewGraphQLServer(authenticationService service.AuthenticationService, noteService service.NoteService, cfg config.Configuration, jwtDatasource domain.JwtDatasource) *Server {
	// Create a new ServeMux
	mux := http.NewServeMux()

	// Create the GraphQL server
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{AuthSrv: authenticationService, NoteSrv: noteService}}))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	// Routes
	if cfg.Environment == "development" {
		mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	}
	mux.Handle("/query", srv)

	var handler http.Handler = mux
	// Add middlewares
	handler = middleware.SetUserInContext(jwtDatasource)(handler)

	// Create the HTTP server
	readTimeOut := 10 * time.Second
	writeTimeOut := 10 * time.Second
	idleTimeOut := 10 * time.Second

	if cfg.Environment == "development" {
		readTimeOut = 10 * time.Minute
		writeTimeOut = 10 * time.Minute
		idleTimeOut = 10 * time.Minute
	}

	httpServer := &http.Server{
		Addr:         net.JoinHostPort("0.0.0.0", cfg.GraphqlServerPort),
		Handler:      handler,
		ReadTimeout:  readTimeOut,
		WriteTimeout: writeTimeOut,
		IdleTimeout:  idleTimeOut,
	}

	return &Server{
		Mux:           mux,
		GraphQLServer: httpServer,
	}
}
