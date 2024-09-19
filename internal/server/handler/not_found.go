package handler

import (
	"fmt"
	"net/http"
)

// NotFoundHandler handles requests to non-existent resources.
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	// Determine the scheme
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	docUrl := fmt.Sprintf("%s://%s/doc", scheme, r.Host)
	msg := fmt.Sprintf("Resource not found. Please refer to the documentation for further details. Doc: %s", docUrl)
	NotFound(w, r, msg)
}
