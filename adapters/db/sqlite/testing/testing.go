package testing

import (
	"context"

	"github.com/jmoiron/sqlx"
	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	_ "modernc.org/sqlite"
)

func NewInMemory() *sqlx.DB {
	db, err := sqlx.Open("sqlite", ":memory:?cache=shared&_pragma=foreign_keys(1)")
	if err != nil {
		panic(err)
	}
	if err := s.ApplyMigrations(context.Background(), db.DB, "up"); err != nil {
		panic(err)
	}
	return db
}
