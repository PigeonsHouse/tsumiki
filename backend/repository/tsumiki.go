package repository

import (
	"database/sql"
	"tsumiki/schema"
)

type TsumikiRepository interface {
	GetTsumiki(watchUserID *int, tsumikiID int) (*schema.Tsumiki, error)
	GetTsumikiBlocks(watchUserID *int, tsumikiID int) ([]schema.TsumikiBlock, error)
	GetTsumikis(watchUserID *int, pageSize, page int, authorID *int, workID *int, keyword string) ([]schema.Tsumiki, error)
	CreateTsumiki(userID int, title string, visibility string, workID *int) (*schema.Tsumiki, error)
	UpdateTsumiki(tsumikiID int, title string, visibility string, workID *int) (*schema.Tsumiki, error)
	DeleteTsumiki(tsumikiID int) error
	CreateMedia(tsumikiID int, mediaType string, url string) (*schema.TsumikiBlockMedia, error)
}

type tsumikiRepositoryImpl struct {
	db *sql.DB
}

func NewTsumikiRepository(db *sql.DB) TsumikiRepository {
	return &tsumikiRepositoryImpl{db: db}
}

func (tr *tsumikiRepositoryImpl) GetTsumiki(watchUserID *int, tsumikiID int) (*schema.Tsumiki, error) {
	return nil, nil
}
func (tr *tsumikiRepositoryImpl) GetTsumikiBlocks(watchUserID *int, tsumikiID int) ([]schema.TsumikiBlock, error) {
	return nil, nil
}
func (tr *tsumikiRepositoryImpl) GetTsumikis(watchUserID *int, pageSize, page int, authorID *int, workID *int, keyword string) ([]schema.Tsumiki, error) {
	return nil, nil
}
func (tr *tsumikiRepositoryImpl) CreateTsumiki(userID int, title string, visibility string, workID *int) (*schema.Tsumiki, error) {
	return nil, nil
}
func (tr *tsumikiRepositoryImpl) UpdateTsumiki(tsumikiID int, title string, visibility string, workID *int) (*schema.Tsumiki, error) {
	return nil, nil
}
func (tr *tsumikiRepositoryImpl) DeleteTsumiki(tsumikiID int) error {
	return nil
}
func (tr *tsumikiRepositoryImpl) CreateMedia(tsumikiID int, mediaType string, url string) (*schema.TsumikiBlockMedia, error) {
	return nil, nil
}
