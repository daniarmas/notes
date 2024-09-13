package service

import (
	"context"
	"errors"

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
	HashDatasource         domain.HashDatasource
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
	// Get the user from the database
	user, err := s.UserRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	// Check if the user password is correct
	if correct := s.HashDatasource.CheckHash(password, user.Password); !correct {
		return nil, errors.New("password incorrect")
	}
	return nil, nil
}
