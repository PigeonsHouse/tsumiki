package repository

import "database/sql"

type DBTX interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

type Repositories struct {
	db                *sql.DB
	User              UserRepository
	Tsumiki           TsumikiRepository
	TsumikiBlock      TsumikiBlockRepository
	TsumikiBlockMedia TsumikiBlockMediaRepository
	Work              WorkRepository
	Thumbnail         ThumbnailRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		db:                db,
		User:              NewUserRepository(db),
		Tsumiki:           NewTsumikiRepository(db),
		TsumikiBlock:      NewTsumikiBlockRepository(db),
		TsumikiBlockMedia: NewTsumikiBlockMediaRepository(db),
		Work:              NewWorkRepository(db),
		Thumbnail:         NewThumbnailRepository(db),
	}
}

func (r *Repositories) withTx(tx *sql.Tx) *Repositories {
	return &Repositories{
		User:              NewUserRepository(tx),
		Tsumiki:           NewTsumikiRepository(tx),
		TsumikiBlock:      NewTsumikiBlockRepository(tx),
		TsumikiBlockMedia: NewTsumikiBlockMediaRepository(tx),
		Work:              NewWorkRepository(tx),
		Thumbnail:         NewThumbnailRepository(tx),
	}
}

type TxCommandFunc = func(txRepos *Repositories) error

func (r *Repositories) RunInTx(fn TxCommandFunc) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	txRepos := r.withTx(tx)

	err = fn(txRepos)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
