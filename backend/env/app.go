package env

import (
	"fmt"
	"os"
	"strconv"
)

var (
	BackendUrl  string
	FrontendUrl string
	AppPort     int
	JwtSecret   []byte
	TLSEnabled  bool
)

func LoadAppEnv() error {
	BackendUrl = os.Getenv("BACKEND_URL")
	if BackendUrl == "" {
		return fmt.Errorf("loading env error: BACKEND_URL")
	}
	FrontendUrl = os.Getenv("FRONTEND_URL")
	if FrontendUrl == "" {
		return fmt.Errorf("loading env error: FRONTEND_URL")
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
	// TLS_ENABLEDが"true"の場合のみTLSを有効化、それ以外はfalse
	TLSEnabled = os.Getenv("TLS_ENABLED") == "true"
	return nil
}
