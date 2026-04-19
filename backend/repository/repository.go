package repository

import "database/sql"

type Repositories struct {
	User              UserRepository
	Tsumiki           TsumikiRepository
	TsumikiBlock      TsumikiBlockRepository
	TsumikiBlockMedia TsumikiBlockMediaRepository
	Work              WorkRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		User:              NewUserRepository(db),
		Tsumiki:           NewTsumikiRepository(db),
		TsumikiBlock:      NewTsumikiBlockRepository(db),
		TsumikiBlockMedia: NewTsumikiBlockMediaRepository(db),
		Work:              NewWorkRepository(db),
	}
}
