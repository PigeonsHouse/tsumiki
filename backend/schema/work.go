package schema

import "time"

type Work struct {
	ID          int              `json:"id"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Visibility  string           `json:"visibility"`
	Thumbnail   *ThumbnailUpload `json:"thumbnail"`
	Owner       User             `json:"owner"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}
