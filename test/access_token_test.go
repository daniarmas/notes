package test

import (
	"testing"

	"github.com/daniarmas/notes/internal/domain"
	"github.com/google/uuid"
)

// Test the NewAccessToken func
func TestNewAccessToken(t *testing.T) {
	t.Run("Test the NewAccessToken func", func(t *testing.T) {
		userId := uuid.New()
		refreshTokenId := uuid.New()
		got := domain.NewAccessToken(userId, refreshTokenId)

		if got.UserId != userId {
			t.Errorf("TestNewAccessToken failed: got UserId %s, want %s", got.UserId, userId)
		}
		if got.RefreshTokenId != refreshTokenId {
			t.Errorf("TestNewAccessToken failed: got RefreshTokenId %s, want %s", got.RefreshTokenId, refreshTokenId)
		}

		// Check if the dynamically generated ID is non-empty
		if got.Id == uuid.Nil {
			t.Errorf("TestNewAccessToken failed: got empty ID")
		}

		// Check if the create and update times are not empty
		if got.CreateTime.IsZero() {
			t.Errorf("TestNewAccessToken failed: got empty CreateTime")
		}
		if got.UpdateTime.IsZero() {
			t.Errorf("TestNewAccessToken failed: got empty UpdateTime")
		}
	})
}
