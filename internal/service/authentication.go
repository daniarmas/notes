package service

import (
	"context"

	"github.com/daniarmas/notes/internal/domain"
)

type SignInResponse struct {
	AccessToken  string
	RefreshToken string
	User         domain.User
}

type AuthenticationService interface {
	SignIn(ctx context.Context, email string, password string) (*SignInResponse, error)
}

type authenticationService struct {
	UserRepository         domain.UserRepository
	AccessTokenRepository  domain.AccessTokenRepository
	RefreshTokenRepository domain.RefreshTokenRepository
}

func NewAuthenticationService(userRepository domain.UserRepository, accessTokenRepository domain.AccessTokenRepository, refreshTokenRepository domain.RefreshTokenRepository) AuthenticationService {
	return &authenticationService{
		UserRepository:         userRepository,
		AccessTokenRepository:  accessTokenRepository,
		RefreshTokenRepository: refreshTokenRepository,
	}
}

func (s *authenticationService) SignIn(ctx context.Context, email string, password string) (*SignInResponse, error) {
	return nil, nil
}
