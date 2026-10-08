package db

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lejeunel/go-image-annotator/config"
	goose "github.com/pressly/goose/v3"
)

type SQLiteDBManager struct {
	*goose.Provider
	*slog.Logger
}

func NewDBManagerFromEnv(logger *slog.Logger) SQLiteDBManager {
	cfg := config.Parse()
	conn := NewSQLiteConnection(DBPath(cfg.LocalArtefactPath))
	manager := NewSQLiteDBManager(conn, logger)
	return manager
}

func NewSQLiteDBManager(db *sqlx.DB, logger *slog.Logger) SQLiteDBManager {
	errCtx := "creating SQLite database manager"
	provider, err := NewMigrationProvider(db.DB)
	if err != nil {
		panic(fmt.Errorf("%v: %w", errCtx, err))
	}

	return SQLiteDBManager{provider, logger}
}

func (m *SQLiteDBManager) Status(ctx context.Context) {
	baseMsg := "SQLite database status"
	statuses, err := m.Provider.Status(ctx)
	if err != nil {
		m.Logger.Warn(baseMsg, "error", err)
		return
	}
	if len(statuses) == 0 {
		m.Logger.Info(baseMsg, "message", "no migrations found")
	}

	for _, status := range statuses {
		appliedAt := "pending"
		if status.State == goose.StateApplied {
			appliedAt = status.AppliedAt.Format(time.RFC3339)
		}
		m.Logger.Info(baseMsg, "applied_at", appliedAt, "source_path", status.Source.Path)
	}
}

func (m *SQLiteDBManager) Down(ctx context.Context) {
	baseMsg := "applying down migration on SQLite database"
	result, err := m.Provider.Down(ctx)
	if errors.Is(err, goose.ErrNoNextVersion) {
		m.Logger.Error(baseMsg, "error", "no applied migrations to roll back")
		return
	}
	if err != nil {
		m.Logger.Error(baseMsg, "error", err.Error())
		return
	}

	m.Logger.Info(baseMsg, "message", result.String())
}

func (m *SQLiteDBManager) Up(ctx context.Context) {
	baseMsg := "applying up migration on SQLite database"
	result, err := m.Provider.UpByOne(ctx)
	if errors.Is(err, goose.ErrNoNextVersion) {
		m.Logger.Error(baseMsg, "error", "no pending migrations to apply")
		return
	}
	if err != nil {
		m.Logger.Error(baseMsg, "error", err.Error())
		return
	}

	m.Logger.Info(baseMsg, "message", result.String())
}

func (m *SQLiteDBManager) Init(ctx context.Context) {
	baseMsg := "initializing SQLite database"
	results, err := m.Provider.Up(ctx)
	if err != nil {
		m.Logger.Error(baseMsg, "error", err.Error())
		panic(err)
	}

	for _, r := range results {
		m.Logger.Info(baseMsg, "message", r.String())
	}
}
