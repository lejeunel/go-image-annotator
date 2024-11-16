package app

import (
	"github.com/jmoiron/sqlx"
	"log"
)

func NewSQLiteConnection(path string) *sqlx.DB {
	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxOpenConns(1)
	return db
}
