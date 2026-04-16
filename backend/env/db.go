package env

import (
	"fmt"
	"os"
)

var (
	MysqlHost     string
	MysqlUser     string
	MysqlPassword string
	MysqlDatabase string
)

func LoadDBEnv() error {
	MysqlHost = os.Getenv("MYSQL_HOST")
	if MysqlHost == "" {
		return fmt.Errorf("loading env error: MYSQL_HOST")
	}
	MysqlUser = os.Getenv("MYSQL_USER")
	if MysqlUser == "" {
		return fmt.Errorf("loading env error: MYSQL_USER")
	}
	MysqlPassword = os.Getenv("MYSQL_PASSWORD")
	if MysqlPassword == "" {
		return fmt.Errorf("loading env error: MYSQL_PASSWORD")
	}
	MysqlDatabase = os.Getenv("MYSQL_DATABASE")
	if MysqlDatabase == "" {
		return fmt.Errorf("loading env error: MYSQL_DATABASE")
	}

	return nil
}
