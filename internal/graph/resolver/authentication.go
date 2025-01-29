package resolver

import (
	"context"
	"errors"

	"github.com/daniarmas/notes/internal/domain"
	"github.com/daniarmas/notes/internal/graph/model"
	"github.com/daniarmas/notes/internal/service"
	"github.com/google/uuid"
)

// map from domain to graphql model user
func mapUser(user domain.User) *model.User {
	var updateTime string
	if !user.UpdateTime.IsZero() {
		updateTime = user.UpdateTime.String()
	}
	return &model.User{
		ID:         user.Id.String(),
		Email:      user.Email,
		Name:       user.Name,
		CreateTime: user.CreateTime.String(),
		UpdateTime: &updateTime,
	}
}

// SignIn is the resolver for the signIn field.
func SignIn(ctx context.Context, input model.SignInInput, srv service.AuthenticationService) (*model.SignInResponse, error) {
	res, err := srv.SignIn(ctx, input.Email, input.Password)
	if err != nil {
		switch err.Error() {
		case "invalid credentials":
			return nil, errors.New("invalid credentials")
		default:
			return nil, errors.New("internal server error")
		}
	}
	return &model.SignInResponse{
		User:         mapUser(res.User),
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

// Me is the resolver for the me field.
func Me(ctx context.Context, srv service.AuthenticationService) (*model.User, error) {
	// Check if the user is authenticated
	userId := domain.GetUserIdFromContext(ctx)
	if userId == uuid.Nil {
		return nil, errors.New("unauthenticated")
	}

	res, err := srv.Me(ctx)
	if err != nil {
		switch err.Error() {
		default:
			return nil, errors.New("internal server error")
		}
	}
	return mapUser(res.User), nil
}
