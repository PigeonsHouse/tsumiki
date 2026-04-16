package store

import "github.com/redis/go-redis/v9"

type Stores struct {
	Auth AuthStore
}

func NewStores(store *redis.Client) *Stores {
	return &Stores{
		Auth: NewAuthStore(store),
	}
}
