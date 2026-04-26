package repository

//go:generate mockgen -source=repository.go -destination=mock/mock_repository.go -package=mock

import "database/sql"

type RowScanner interface {
	Scan(dest ...any) error
}

type RowsScanner interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
	Err() error
}

type DBTX interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (RowsScanner, error)
	QueryRow(query string, args ...any) RowScanner
}


type Repositories struct {
	db                *sql.DB
	RunTxFn           func(fn TxCommandFunc) error // nil = use real db transaction; override in tests
	User              UserRepository
	Tsumiki           TsumikiRepository
	TsumikiBlock      TsumikiBlockRepository
	TsumikiBlockMedia TsumikiBlockMediaRepository
	Work              WorkRepository
	Thumbnail         ThumbnailRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	adapted := &dbTXAdapter{inner: db}
	return &Repositories{
		db:                db,
		User:              NewUserRepository(adapted),
		Tsumiki:           NewTsumikiRepository(adapted),
		TsumikiBlock:      NewTsumikiBlockRepository(adapted),
		TsumikiBlockMedia: NewTsumikiBlockMediaRepository(adapted),
		Work:              NewWorkRepository(adapted),
		Thumbnail:         NewThumbnailRepository(adapted),
	}
}

func (r *Repositories) withTx(tx *sql.Tx) *Repositories {
	adapted := &dbTXAdapter{inner: tx}
	return &Repositories{
		User:              NewUserRepository(adapted),
		Tsumiki:           NewTsumikiRepository(adapted),
		TsumikiBlock:      NewTsumikiBlockRepository(adapted),
		TsumikiBlockMedia: NewTsumikiBlockMediaRepository(adapted),
		Work:              NewWorkRepository(adapted),
		Thumbnail:         NewThumbnailRepository(adapted),
	}
}

type TxCommandFunc = func(txRepos *Repositories) error

func (r *Repositories) RunInTx(fn TxCommandFunc) error {
	if r.RunTxFn != nil {
		return r.RunTxFn(fn)
	}
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
