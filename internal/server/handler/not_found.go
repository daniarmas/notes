package handler

import "net/http"

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	NotFound(w, r, "Please check the endpoint or resource ID and try again.")
}
