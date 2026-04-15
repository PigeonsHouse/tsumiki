package schema

import "time"

type User struct {
	ID            string    `db:"id" json:"id"`
	DiscordUserID int64     `db:"discord_user_id" json:"-"`
	Name          string    `db:"name" json:"name"`
	ThumbnailUrl  string    `db:"thumbnail_url" json:"thumbnail_url"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}
