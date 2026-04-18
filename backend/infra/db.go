package infra

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"fmt"
	"tsumiki/env"
)

func NewDatabase() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", env.MysqlUser, env.MysqlPassword, env.MysqlHost, env.MysqlDatabase)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
