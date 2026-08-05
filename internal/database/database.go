package database

import (
	"database/sql"
	"g0/internal/config"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Connect() error {
	var err error

	DB, err = sql.Open("sqlite", config.DatabaseUrl)
	if err != nil {
		return err
	}

	DB.SetMaxOpenConns(1)
	DB.SetMaxIdleConns(1)

	return DB.Ping()
}
