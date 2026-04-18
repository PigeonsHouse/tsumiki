package repository

import "database/sql"

type Repositories struct {
	Auth AuthRepository
	User UserRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Auth: NewAuthRepository(db),
		User: NewUserRepository(db),
	}
}
