package redis

import (
	"context"
	"testing"
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
)

func TestTokenRepositoryCreateRejectsExpiredSession(t *testing.T) {
	repository := TokenRepository{}

	err := repository.Create(context.Background(), entity.TokenSession{
		AccessExpiresAt:  time.Now().Add(-time.Second),
		RefreshExpiresAt: time.Now().Add(time.Minute),
	})
	if err == nil {
		t.Fatal("expected an expired session to be rejected")
	}
}
