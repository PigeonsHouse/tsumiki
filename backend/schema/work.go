package schema

import "time"

type Work struct {
	ID           int       `db:"id" json:"id"`
	Title        string    `db:"title" json:"title"`
	Description  string    `db:"description" json:"description"`
	ThumbnailUrl *string   `db:"thumbnail_url" json:"thumbnail_url"`
	Owner        User      `json:"owner"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
