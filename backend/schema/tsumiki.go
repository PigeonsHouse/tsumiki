package schema

import "time"

type Tsumiki struct {
	ID           int       `db:"id" json:"id"`
	Title        string    `db:"title" json:"title"`
	ThumbnailUrl *string   `db:"thumbnail_url" json:"thumbnail_url"`
	Visibility   string    `db:"visibility" json:"visibility"`
	User         User      `json:"user"`
	Work         *Work     `json:"work"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type TsumikiBlock struct {
	ID          int               `db:"id" json:"id"`
	Message     *string           `db:"message" json:"message"`
	Medias      TsumikiBlockMedia `json:"medias"`
	Percentage  int               `db:"percentage" json:"percentage"`
	Condition   int               `db:"condition" json:"condition"`
	NextBlockId *int              `db:"next_block_id" json:"-"`
	TsumikiId   int               `db:"tsumiki_id" json:"-"`
	CreatedAt   time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time         `db:"updated_at" json:"updated_at"`
}

type TsumikiBlockMedia struct {
	ID        int       `db:"id" json:"id"`
	Type      string    `db:"type" json:"type"`
	Url       string    `db:"url" json:"url"`
	Order     int       `db:"order" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
