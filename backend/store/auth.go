package store

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const refreshTokenTTL = 30 * 24 * time.Hour

type AuthStore interface {
	SetRefreshToken(ctx context.Context, userID int, sessionID string) error
}

type authStoreImpl struct {
	store *redis.Client
}

func NewAuthStore(store *redis.Client) AuthStore {
	return &authStoreImpl{
		store: store,
	}
}

func (as *authStoreImpl) SetRefreshToken(ctx context.Context, userID int, sessionID string) error {
	key := fmt.Sprintf("refresh_token:%d:%s", userID, sessionID)
	return as.store.Set(ctx, key, true, refreshTokenTTL).Err()
}
