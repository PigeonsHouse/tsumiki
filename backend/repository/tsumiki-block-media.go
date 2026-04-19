package repository

import (
	"database/sql"
	"tsumiki/schema"
)

type TsumikiBlockMediaRepository interface {
	CreateMedia(url string, mediaType string) (*schema.TsumikiBlockMedia, error)
	SetMediaRelation(blockID int, updatedMediaIDs []int) error
}

type tsumikiBlockMediaRepositoryImpl struct {
	db *sql.DB
}

func NewTsumikiBlockMediaRepository(db *sql.DB) TsumikiBlockMediaRepository {
	return &tsumikiBlockMediaRepositoryImpl{db: db}
}

func (tbmr *tsumikiBlockMediaRepositoryImpl) CreateMedia(url string, mediaType string) (*schema.TsumikiBlockMedia, error) {
	return nil, nil
}
func (tbmr *tsumikiBlockMediaRepositoryImpl) SetMediaRelation(blockID int, updatedMediaIDs []int) error {
	return nil
}
