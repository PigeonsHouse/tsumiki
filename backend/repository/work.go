package repository

import (
	"database/sql"
)

type WorkRepository interface {
}

type workRepositoryImpl struct {
	db *sql.DB
}

func NewWorkRepository(db *sql.DB) WorkRepository {
	return &workRepositoryImpl{
		db: db,
	}
}
