package db

import (
	"context"
	"database/sql"
	"embed"
	"github.com/jmoiron/sqlx"
	goose "github.com/pressly/goose/v3"
	"io/fs"
	_ "modernc.org/sqlite"
	"net/url"
	"os"
	"path/filepath"
)

//go:embed migrations/*.sql
var MigrationsFS embed.FS

func NewSQLiteConnection(path string) *sqlx.DB {
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		panic(err)
	}
	q := url.Values{
		"_time_format": {"sqlite"},
		"_pragma": {
			"foreign_keys(ON)",
			"journal_mode(WAL)",
			"synchronous(NORMAL)",
			"busy_timeout(5000)",
			"journal_size_limit(1000000)",
			"mmap_size(30000000000)",
			"cache_size(-2000)",
		},
	}

	u := url.URL{
		Scheme:   "file",
		Path:     path,
		RawQuery: q.Encode(),
	}

	db, err := sqlx.Open("sqlite", u.String())
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(4)
	db.SetMaxIdleConns(4)

	return db
}

func NewMigrationProvider(db *sql.DB) (*goose.Provider, error) {

	migrationsFSsub, err := fs.Sub(MigrationsFS, "migrations")
	if err != nil {
		return nil, err
	}
	provider, err := goose.NewProvider(goose.DialectSQLite3, db, migrationsFSsub)
	if err != nil {
		return nil, err
	}
	return provider, nil

}

func ApplyMigrations(ctx context.Context, db *sql.DB, direction string) error {
	provider, err := NewMigrationProvider(db)
	if err != nil {
		panic(err)
	}
	switch direction {
	case "up":
		_, err := provider.Up(ctx)
		if err != nil {
			return err
		}

	case "down":
		_, err := provider.Down(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewInMemory() *sqlx.DB {
	db, err := sqlx.Open("sqlite", ":memory:?cache=shared&_pragma=foreign_keys(1)")
	if err != nil {
		panic(err)
	}
	if err := ApplyMigrations(context.Background(), db.DB, "up"); err != nil {
		panic(err)
	}
	return db
}

func NewSQLiteDB(path string) *sqlx.DB {
	db := NewSQLiteConnection(path)
	if err := ApplyMigrations(context.Background(), db.DB, "up"); err != nil {
		panic(err)
	}
	return db
}
