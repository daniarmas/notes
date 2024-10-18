package utils

import (
	"context"
	"path/filepath"
	"runtime"

)

type contextKey string

const RequestIDKey contextKey = "request-id"

// GetFileName returns the file name of the caller
func GetFileName() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return "unknown"
	}
	return filepath.Base(file)
}

// GetFunctionName returns the function name of the caller
func GetFunctionName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "unknown"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}

	return filepath.Base(fn.Name())
}

// GetLineNumber returns the line number of the caller
func GetLineNumber() int {
	_, _, line, ok := runtime.Caller(1)
	if !ok {
		return 0
	}
	return line
}

func ExtractRequestIdFromContext(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return "unknown"

}
