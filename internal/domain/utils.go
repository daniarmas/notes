package domain

import "context"

// SetUserInContext sets the user in the context
func SetUserInContext(ctx context.Context, user *User) context.Context {
	if user == nil {
		return ctx
	} else {
		return context.WithValue(ctx, "user", user)
	}
}

// Get user from context
func GetUserFromContext(ctx context.Context) *User {
	if user, ok := ctx.Value("user").(User); ok {
		return &user
	}
	return nil
}
