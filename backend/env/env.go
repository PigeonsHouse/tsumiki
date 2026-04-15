package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	BackendUrl  string
	FrontendUrl string

	AppPort       int
	JwtSecret     string
	AllowGuildIds []string

	DiscordClientID string
	DiscordSecretID string
)

func LoadEnv() error {
	BackendUrl = os.Getenv("BACKEND_URL")
	if BackendUrl == "" {
		return fmt.Errorf("loading env error: BACKEND_URL")
	}
	FrontendUrl = os.Getenv("FRONTEND_URL")
	if FrontendUrl == "" {
		FrontendUrl = BackendUrl
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

	DiscordClientID = os.Getenv("DISCORD_CLIENT_ID")
	if DiscordClientID == "" {
		return fmt.Errorf("loading env error: DISCORD_CLIENT_ID")
	}
	DiscordSecretID = os.Getenv("DISCORD_SECRET_ID")
	if DiscordSecretID == "" {
		return fmt.Errorf("loading env error: DISCORD_SECRET_ID")
	}

	return nil
}
