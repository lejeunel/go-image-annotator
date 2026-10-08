package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
	goose "github.com/pressly/goose/v3"

	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var MigrationsFS embed.FS

// DBFileName is the name of the sqlite database file held in the local
// artefact directory.
const DBFileName = "db.sqlite"

// DBPath returns the path of the database file within the local artefact
// directory localPath.
func DBPath(localPath string) string {
	return filepath.Join(localPath, DBFileName)
}

func NewSQLiteConnection(path string) *sqlx.DB {
	if path == "" {
		panic("sqlite: database path is empty")
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		panic(fmt.Sprintf("sqlite: cannot resolve %q: %v", path, err))
	}

	if info, err := os.Stat(abs); err == nil && info.IsDir() {
		panic(fmt.Sprintf("sqlite: %q is a directory, expected a database file", abs))
	}

	dir := filepath.Dir(abs)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		panic(fmt.Sprintf("sqlite: cannot create directory %q for database %q: %v", dir, path, err))
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

	// Absolute + forward slashes + leading "/" gives file:///abs/path on Unix
	// and file:///C:/abs/path on Windows.
	uriPath := filepath.ToSlash(abs)
	if !strings.HasPrefix(uriPath, "/") {
		uriPath = "/" + uriPath
	}
	u := url.URL{
		Scheme:   "file",
		Path:     uriPath,
		RawQuery: q.Encode(),
	}

	db, err := sqlx.Open("sqlite", u.String())
	if err != nil {
		panic(fmt.Sprintf("sqlite: open %q: %v", abs, err))
	}
	if err := db.Ping(); err != nil {
		db.Close()
		panic(fmt.Sprintf("sqlite: cannot connect to %q (from %q): %v", abs, path, err))
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
