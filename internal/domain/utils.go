package domain

import (
	"context"

	"github.com/google/uuid"
)

// SetUserInContext sets the user in the context
func SetUserInContext(ctx context.Context, userId uuid.UUID) context.Context {
	if userId == uuid.Nil {
		return ctx
	} else {
		return context.WithValue(ctx, "userId", userId.String())
	}
}

// Get user from context
func GetUserIdFromContext(ctx context.Context) uuid.UUID {
	if user, ok := ctx.Value("userId").(string); ok {
		return uuid.MustParse(user)
	}
	return uuid.Nil
}
