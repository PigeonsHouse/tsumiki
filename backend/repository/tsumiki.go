package repository

import (
	"database/sql"
)

type TsumikiRepository interface {
}

type tsumikiRepositoryImpl struct {
	db *sql.DB
}

func NewTsumikiRepository(db *sql.DB) TsumikiRepository {
	return &tsumikiRepositoryImpl{
		db: db,
	}
}
