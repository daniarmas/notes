package service

import (
	"context"
	"errors"
	"time"

	"github.com/daniarmas/notes/internal/customerrors"
	"github.com/daniarmas/notes/internal/domain"
)

type SignInResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	User         domain.User `json:"user"`
}

type AuthenticationService interface {
	SignIn(ctx context.Context, email string, password string) (*SignInResponse, error)
	SignOut(ctx context.Context) error
}

type authenticationService struct {
	JwtDatasource          domain.JwtDatasource
	HashDatasource         domain.HashDatasource
	UserRepository         domain.UserRepository
	AccessTokenRepository  domain.AccessTokenRepository
	RefreshTokenRepository domain.RefreshTokenRepository
}

func NewAuthenticationService(jwtDatasource domain.JwtDatasource, hashDatasource domain.HashDatasource, userRepository domain.UserRepository, accessTokenRepository domain.AccessTokenRepository, refreshTokenRepository domain.RefreshTokenRepository) AuthenticationService {
	return &authenticationService{
		UserRepository:         userRepository,
		AccessTokenRepository:  accessTokenRepository,
		RefreshTokenRepository: refreshTokenRepository,
		HashDatasource:         hashDatasource,
		JwtDatasource:          jwtDatasource,
	}
}

func (s *authenticationService) SignIn(ctx context.Context, email string, password string) (*SignInResponse, error) {
	// Get the user by email
	user, err := s.UserRepository.GetUserByEmail(ctx, email)
	if err != nil {
		switch err.(type) {
		case *customerrors.RecordNotFound:
			return nil, errors.New("invalid credentials")
		default:
			return nil, err
		}
	}
	// Check if the user password is correct
	correct, err := s.HashDatasource.CheckHash(password, user.Password)
	if err != nil {
		return nil, err
	}

	// If the password is incorrect, return an error
	if !correct {
		return nil, errors.New("invalid credentials")
	}

	// Delete the existing access token
	err = s.AccessTokenRepository.DeleteAccessTokenByUserId(ctx, user.Id)
	if err != nil {
		switch err.(type) {
		case *customerrors.RecordNotFound:
			// Do nothing
		default:
			return nil, err
		}
	}
	// Delete the existing refresh token
	err = s.RefreshTokenRepository.DeleteRefreshTokenByUserId(ctx, user.Id)
	if err != nil {
		switch err.(type) {
		case *customerrors.RecordNotFound:
			// Do nothing
		default:
			return nil, err
		}
	}
	// Create a new refresh token
	refreshToken, err := s.RefreshTokenRepository.CreateRefreshToken(ctx, &domain.RefreshToken{
		UserId: user.Id,
	})
	if err != nil {
		return nil, err
	}
	// Create a new access token
	accessToken, err := s.AccessTokenRepository.CreateAccessToken(ctx, user.Id, refreshToken.Id)
	if err != nil {
		return nil, err
	}
	// Create jwt
	now := time.Now()
	accessTokenExpiration := now.Add(60 * time.Minute)
	refreshTokenExpiration := now.Add(30 * 24 * time.Hour)
	// Refresh token jwt
	refreshTokenJWT, err := s.JwtDatasource.CreateJWT(&domain.JWTMetadata{TokenId: refreshToken.Id}, refreshTokenExpiration)
	if err != nil {
		return nil, err
	}
	// Refresh token jwt
	accessTokenJWT, err := s.JwtDatasource.CreateJWT(&domain.JWTMetadata{TokenId: accessToken.Id}, accessTokenExpiration)
	if err != nil {
		return nil, err
	}
	return &SignInResponse{
		AccessToken:  *accessTokenJWT,
		RefreshToken: *refreshTokenJWT,
		User:         *user,
	}, nil
}

func (s *authenticationService) SignOut(ctx context.Context) error {
	// Get the user from the context
	user := domain.GetUserFromContext(ctx)
	// Delete the existing access token
	err := s.AccessTokenRepository.DeleteAccessTokenByUserId(ctx, user.Id)
	if err != nil {
		return err
	}
	// Delete the existing refresh token
	err = s.RefreshTokenRepository.DeleteRefreshTokenByUserId(ctx, user.Id)
	if err != nil {
		return err
	}
	return nil
}
