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
	ValidateAndDeleteRefreshToken(ctx context.Context, userID int, sessionID string) (bool, error)
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

func (as *authStoreImpl) ValidateAndDeleteRefreshToken(ctx context.Context, userID int, sessionID string) (bool, error) {
	key := fmt.Sprintf("refresh_token:%d:%s", userID, sessionID)
	err := as.store.GetDel(ctx, key).Err()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
