package repository

import (
	"database/sql"
	"tsumiki/schema"
)

type TsumikiBlockRepository interface {
	IsBelongToTsumiki(tsumikiID int, blockID int) (bool, error)
	CreateBlock(tsumikiID int, message *string, percentage int, condition int, mediaIDs []int) (*schema.TsumikiBlock, error)
	UpdateBlock(blockID int, message *string, percentage int, condition int, mediaIDs []int) (*schema.TsumikiBlock, error)
	SoftDeleteBlock(blockID int) error
}

type tsumikiBlockRepositoryImpl struct {
	db *sql.DB
}

func NewTsumikiBlockRepository(db *sql.DB) TsumikiBlockRepository {
	return &tsumikiBlockRepositoryImpl{db: db}
}

func (tbr *tsumikiBlockRepositoryImpl) IsBelongToTsumiki(tsumikiID int, blockID int) (bool, error) {
	return false, nil
}
func (tbr *tsumikiBlockRepositoryImpl) CreateBlock(tsumikiID int, message *string, percentage int, condition int, mediaIDs []int) (*schema.TsumikiBlock, error) {
	return nil, nil
}
func (tbr *tsumikiBlockRepositoryImpl) UpdateBlock(blockID int, message *string, percentage int, condition int, mediaIDs []int) (*schema.TsumikiBlock, error) {
	return nil, nil
}
func (tbr *tsumikiBlockRepositoryImpl) SoftDeleteBlock(blockID int) error { return nil }
