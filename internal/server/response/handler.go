package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message any         `json:"message"`
	Details any         `json:"details"`
	Data    interface{} `json:"data"`
}

func UnauthorizedHandler(w http.ResponseWriter, r *http.Request, message string, errors map[string]string) {
	res := Response{
		Code:    http.StatusUnauthorized,
		Message: message,
		Details: errors,
		Data:    nil,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(res)
}

func BadRequestHandler(w http.ResponseWriter, r *http.Request, message string, errors map[string]string) {
	res := Response{
		Code:    http.StatusBadRequest,
		Message: message,
		Details: errors,
		Data:    nil,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(res)
}

func StatusOk(w http.ResponseWriter, r *http.Request, data any) {
	res := Response{
		Code:    http.StatusOK,
		Message: "OK",
		Details: nil,
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	res := Response{
		Code:    http.StatusInternalServerError,
		Message: "Internal Server Error",
		Details: nil,
		Data:    nil,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(res)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	res := Response{
		Code:    http.StatusNotFound,
		Message: "Not Found",
		Details: nil,
		Data:    nil,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(res)
}
