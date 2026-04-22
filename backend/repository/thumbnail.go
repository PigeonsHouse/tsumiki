package repository

import (
	"database/sql"
	"tsumiki/schema"
)

type ThumbnailRepository interface {
	Create(userID int, path string) (*schema.ThumbnailUpload, error)
	Get(thumbnailID int) (*schema.ThumbnailUpload, error)
	IsInUse(thumbnailID int) (bool, error)
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

func (r *thumbnailRepositoryImpl) IsInUse(thumbnailID int) (bool, error) {
	var exists int
	err := r.db.QueryRow(
		"SELECT 1 FROM tsumikis WHERE thumbnail_id = ? UNION SELECT 1 FROM works WHERE thumbnail_id = ? LIMIT 1",
		thumbnailID, thumbnailID,
	).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
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
