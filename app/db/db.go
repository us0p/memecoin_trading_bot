package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	db *sql.DB
}

func NewDB(options string) (DB, error) {
	dbCon, err := sql.Open("sqlite3", options)
	if err != nil {
		return DB{}, err
	}

	db := DB{
		dbCon,
	}

	return db, nil
}
