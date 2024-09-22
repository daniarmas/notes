package handler

import (
	"fmt"
	"net/http"
	"strings"
)

// NotFoundHandler handles requests to non-existent resources.
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	// Determine the scheme
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	docUrl := fmt.Sprintf("%s://%s/doc", scheme, r.Host)

	acceptHeader := r.Header.Get("Accept")
	if acceptHeader != "" && (strings.Contains(acceptHeader, "text/html") || strings.Contains(acceptHeader, "application/xhtml+xml")) {
		http.Redirect(w, r, docUrl, http.StatusFound)
		return
	}

	msg := fmt.Sprintf("Resource not found. Please refer to the documentation for further details. Doc: %s", docUrl)
	NotFound(w, r, msg)
}
