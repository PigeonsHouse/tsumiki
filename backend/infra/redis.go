package infra

import (
	"context"
	"fmt"
	"tsumiki/env"

	"github.com/redis/go-redis/v9"
)

func NewRedis() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", env.RedisHost),
		Password: env.RedisPassword,
		DB:       0,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	return rdb, nil
}
