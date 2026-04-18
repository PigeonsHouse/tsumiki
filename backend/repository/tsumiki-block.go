package repository

import (
	"database/sql"
)

type TsumikiBlockRepository interface {
}

type tsumikiBlockRepositoryImpl struct {
	db *sql.DB
}

func NewTsumikiBlockRepository(db *sql.DB) TsumikiBlockRepository {
	return &tsumikiBlockRepositoryImpl{
		db: db,
	}
}
