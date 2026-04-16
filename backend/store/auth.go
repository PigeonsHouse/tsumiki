package store

import "github.com/redis/go-redis/v9"

type AuthStore interface {
	SetRefreshToken()
}

type AuthStoreImpl struct {
	store *redis.Client
}

func NewAuthStore(store *redis.Client) AuthStore {
	return &AuthStoreImpl{
		store: store,
	}
}

func (as *AuthStoreImpl) SetRefreshToken() {}
