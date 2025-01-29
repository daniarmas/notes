package graph

import (
	"github.com/daniarmas/notes/internal/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	AuthSrv service.AuthenticationService
}
