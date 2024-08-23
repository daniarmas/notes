package server

import (
	"net/http"

	"github.com/daniarmas/notes/internal/config"
)

// Instantiate a net/http server...
func NewServer(
	config *config.Configuration,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		config,
	)
	var handler http.Handler = mux
	return handler
}

func addRoutes(
	mux *http.ServeMux,
	config *config.Configuration,
) {
	mux.Handle("/", http.NotFoundHandler())
}
