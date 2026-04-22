package repository

import (
	"database/sql"
	"tsumiki/schema"
)

type ThumbnailRepository interface {
	Create(userID int, path string) (*schema.ThumbnailUpload, error)
	Get(thumbnailID int) (*schema.ThumbnailUpload, error)
}

type thumbnailRepositoryImpl struct {
	db DBTX
}

func NewThumbnailRepository(db DBTX) ThumbnailRepository {
	return &thumbnailRepositoryImpl{db: db}
}

func (r *thumbnailRepositoryImpl) Create(userID int, path string) (*schema.ThumbnailUpload, error) {
	result, err := r.db.Exec(
		"INSERT INTO thumbnails (user_id, path) VALUES (?, ?)",
		userID, path,
	)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.Get(int(id))
}

func (r *thumbnailRepositoryImpl) Get(thumbnailID int) (*schema.ThumbnailUpload, error) {
	var t schema.ThumbnailUpload
	err := r.db.QueryRow(
		"SELECT id, path, created_at, updated_at FROM thumbnails WHERE id = ?",
		thumbnailID,
	).Scan(&t.ID, &t.Url, &t.CreatedAt, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}
