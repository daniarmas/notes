package service

import "github.com/daniarmas/notes/internal/domain"

type SignInResponse struct {
	AccessToken  string
	RefreshToken string
	User         domain.User
}

type AuthenticationService interface {
	SignIn(email string, password string) (*SignInResponse, error)
}
