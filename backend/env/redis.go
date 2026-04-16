package env

import (
	"fmt"
	"os"
)

var (
	RedisHost     string
	RedisPassword string
)

func LoadRedisEnv() error {
	RedisHost = os.Getenv("REDIS_HOST")
	if RedisHost == "" {
		return fmt.Errorf("loading env error: REDIS_HOST")
	}
	RedisPassword = os.Getenv("REDIS_PASSWORD")
	if RedisPassword == "" {
		return fmt.Errorf("loading env error: REDIS_PASSWORD")
	}

	return nil
}
