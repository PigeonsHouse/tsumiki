package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	BackendUrl    string
	AppPort       int
	JwtSecret     string
	AllowGuildIds []string
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
	JwtSecret = os.Getenv("JWT_SECRET")
	if JwtSecret == "" {
		return fmt.Errorf("loading env error: JWT_SECRET")
	}
	AllowGuildIdsStr := os.Getenv("ALLOW_GUILD_IDS")
	AllowGuildIds = strings.Split(AllowGuildIdsStr, ",")
	return nil
}
