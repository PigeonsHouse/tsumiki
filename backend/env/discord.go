package env

import (
	"fmt"
	"os"
)

var (
	DiscordClientID string
	DiscordSecretID string
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

	return nil
}
