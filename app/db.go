package app

import (
	"database/sql"
	c "datahub/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"log/slog"
	_ "modernc.org/sqlite"
	"os"
	"time"
)

func setPragma(db *sql.DB, pragma, value string, logger *slog.Logger) {
	_, err := db.Exec("PRAGMA " + pragma + "=" + value)
	if err != nil {
		logger.Error("Failed to set PRAGMA", "option",
			pragma, "value", value, "error", err)
		os.Exit(1)
	}
}

func NewSQLiteConnection(path string, logger *slog.Logger) *sqlx.DB {
	logger.Info("Opening SQLite db", "path", path)
	db, err := sqlx.Open("sqlite", path)
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	setPragma(db.DB, "foreign_keys", "ON", logger)
	setPragma(db.DB, "journal_mode", "WAL", logger)
	setPragma(db.DB, "synchronous", "NORMAL", logger)
	setPragma(db.DB, "busy_timeout", "5000", logger)
	setPragma(db.DB, "journal_size_limit", "1000000", logger)
	setPragma(db.DB, "mmap_size", "30000000000", logger)
	setPragma(db.DB, "cache_size", "-2000", logger)
	return db
}

func NewPostgresConnection(cfg *c.Config, logger *slog.Logger) *sqlx.DB {

	logger.Info("Opening connection to PostgreSQL database",
		"host", cfg.PostGreSqlHost,
		"port", cfg.PostGreSqlPort)
	db, err := sqlx.Open("pgx", c.NewConfig().DBDataSourceName())
	if err != nil {
		log.Fatalln(err)
	}
	return db
}
