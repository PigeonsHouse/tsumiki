package repository

import "database/sql"

type Repositories struct {
	Auth AuthRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Auth: NewAuthRepository(db),
	}
}
