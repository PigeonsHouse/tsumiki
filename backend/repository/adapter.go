package repository

import "database/sql"

// sqlDBTX は *sql.DB / *sql.Tx が共通して持つメソッドセット。
// 標準ライブラリの具体型を DBTX に適合させるためのブリッジとして使う。
type sqlDBTX interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

type dbTXAdapter struct {
	inner sqlDBTX
}

func (a *dbTXAdapter) Exec(query string, args ...any) (sql.Result, error) {
	return a.inner.Exec(query, args...)
}

func (a *dbTXAdapter) Query(query string, args ...any) (RowsScanner, error) {
	return a.inner.Query(query, args...)
}

func (a *dbTXAdapter) QueryRow(query string, args ...any) RowScanner {
	return a.inner.QueryRow(query, args...)
}
