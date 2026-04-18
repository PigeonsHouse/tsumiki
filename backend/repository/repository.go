package repository

import "database/sql"

type Repositories struct {
	User    UserRepository
	Tsumiki TsumikiRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		User:    NewUserRepository(db),
		Tsumiki: NewTsumikiRepository(db),
	}
}
