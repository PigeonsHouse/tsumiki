package schema

import "time"

type User struct {
	ID            int       `db:"id" json:"id"`
	DiscordUserID string    `db:"discord_user_id" json:"-"`
	Name          string    `db:"name" json:"name"`
	GuildID       *string   `db:"guild_id" json:"-"`
	AvatarUrl     string    `db:"avatar_url" json:"avatar_url"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}
