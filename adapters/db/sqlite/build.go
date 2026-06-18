package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"github.com/jmoiron/sqlx"
	goose "github.com/pressly/goose/v3"
	"io/fs"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"time"
)

//go:embed migrations/*.sql
var MigrationsFS embed.FS

func NewSQLiteConnection(path string) *sqlx.DB {
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		panic(err)
	}
	db, err := sqlx.Open("sqlite", fmt.Sprintf("file:%s?_time_format=sqlite", path))
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	setPragma(db.DB, "foreign_keys", "ON")
	setPragma(db.DB, "journal_mode", "WAL")
	setPragma(db.DB, "synchronous", "NORMAL")
	setPragma(db.DB, "busy_timeout", "5000")
	setPragma(db.DB, "journal_size_limit", "1000000")
	setPragma(db.DB, "mmap_size", "30000000000")
	setPragma(db.DB, "cache_size", "-2000")
	return db
}

func setPragma(db *sql.DB, pragma, value string) {
	_, err := db.Exec("PRAGMA " + pragma + "=" + value)
	if err != nil {
		panic(err)
	}
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

func ApplyMigrations(ctx context.Context, provider *goose.Provider, direction string) error {
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

func NewSQLiteDB(path string) *sqlx.DB {
	db := NewSQLiteConnection(path)
	m, err := NewMigrationProvider(db.DB)
	if err != nil {
		panic(err)
	}
	if err := ApplyMigrations(context.Background(), m, "up"); err != nil {
		panic(err)
	}
	return db

}
