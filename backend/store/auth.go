package store

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const refreshTokenTTL = 30 * 24 * time.Hour

type AuthStore interface {
	SetRefreshToken(ctx context.Context, userID, sessionID string) error
}

type AuthStoreImpl struct {
	store *redis.Client
}

func NewAuthStore(store *redis.Client) AuthStore {
	return &AuthStoreImpl{
		store: store,
	}
}

func (as *AuthStoreImpl) SetRefreshToken(ctx context.Context, userID, sessionID string) error {
	key := fmt.Sprintf("refresh_token:%s:%s", userID, sessionID)
	return as.store.Set(ctx, key, true, refreshTokenTTL).Err()
}
