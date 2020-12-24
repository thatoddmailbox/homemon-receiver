package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func connectToDatabase() error {
	var err error
	db, err = sql.Open("mysql", currentConfig.DatabaseDSN)
	if err != nil {
		return err
	}

	return db.Ping()
}
