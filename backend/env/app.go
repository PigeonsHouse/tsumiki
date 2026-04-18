package env

import (
	"fmt"
	"os"
	"strconv"
)

var (
	BackendUrl string
	AppPort    int
	JwtSecret  []byte
)

func LoadAppEnv() error {
	BackendUrl = os.Getenv("BACKEND_URL")
	if BackendUrl == "" {
		return fmt.Errorf("loading env error: BACKEND_URL")
	}
	if port := os.Getenv("APP_PORT"); port != "" {
		portNum, err := strconv.Atoi(port)
		if err != nil {
			return fmt.Errorf("loading env error: APP_PORT: %w", err)
		}
		AppPort = portNum
	} else {
		AppPort = 8000
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return fmt.Errorf("loading env error: JWT_SECRET")
	}
	JwtSecret = []byte(jwtSecret)
	return nil
}
