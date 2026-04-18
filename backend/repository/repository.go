package repository

import "database/sql"

type Repositories struct {
	User UserRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		User: NewUserRepository(db),
	}
}
