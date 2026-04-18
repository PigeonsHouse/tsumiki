package env

import (
	"fmt"
	"os"
	"strings"
)

var (
	DiscordClientID string
	DiscordSecretID string

	AllowGuildIds []string
)

func LoadDiscordEnv() error {
	DiscordClientID = os.Getenv("DISCORD_CLIENT_ID")
	if DiscordClientID == "" {
		return fmt.Errorf("loading env error: DISCORD_CLIENT_ID")
	}
	DiscordSecretID = os.Getenv("DISCORD_SECRET_ID")
	if DiscordSecretID == "" {
		return fmt.Errorf("loading env error: DISCORD_SECRET_ID")
	}
	AllowGuildIdsStr := os.Getenv("ALLOW_GUILD_IDS")
	AllowGuildIds = strings.Split(AllowGuildIdsStr, ",")

	return nil
}
