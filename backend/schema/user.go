package schema

import "time"

type User struct {
	ID            int       `db:"id" json:"id"`
	DiscordUserID string    `db:"discord_user_id" json:"-"`
	Name          string    `db:"name" json:"name"`
	ThumbnailUrl  string    `db:"thumbnail_url" json:"thumbnail_url"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}
