package schema

import "time"

type ThumbnailUpload struct {
	ID        int       `json:"id"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Tsumiki struct {
	ID         int              `json:"id"`
	Title      string           `json:"title"`
	Thumbnail  *ThumbnailUpload `json:"thumbnail"`
	Visibility string           `json:"visibility"`
	User       User             `json:"user"`
	Work       *Work            `json:"work"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}

type TsumikiBlock struct {
	ID          int                 `db:"id" json:"id"`
	Message     *string             `db:"message" json:"message"`
	Medias      []TsumikiBlockMedia `json:"medias"`
	Percentage  int                 `db:"percentage" json:"percentage"`
	Condition   int                 `db:"condition" json:"condition"`
	NextBlockId *int                `db:"next_block_id" json:"-"`
	TsumikiId   int                 `db:"tsumiki_id" json:"-"`
	CreatedAt   time.Time           `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time           `db:"updated_at" json:"updated_at"`
}

// ブロック一覧取得用。削除済みブロックは is_deleted のみ返す
type TsumikiBlockView struct {
	ID         int                 `json:"id"`
	IsDeleted  bool                `json:"is_deleted"`
	Message    *string             `json:"message,omitempty"`
	Medias     []TsumikiBlockMedia `json:"medias,omitempty"`
	Percentage *int                `json:"percentage,omitempty"`
	Condition  *int                `json:"condition,omitempty"`
	CreatedAt  *time.Time          `json:"created_at,omitempty"`
	UpdatedAt  *time.Time          `json:"updated_at,omitempty"`
}

type TsumikiBlockMedia struct {
	ID        int       `db:"id" json:"id"`
	Type      string    `db:"type" json:"type"`
	Url       string    `db:"url" json:"url"`
	Order     int       `db:"order" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
