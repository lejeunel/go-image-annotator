package migrations

import (
	"context"
	goose "github.com/pressly/goose/v3"
	"log/slog"
)

func LogMigrationResult(r *goose.MigrationResult, logger *slog.Logger) {
	logger.Info("Applying migration", "result", r.String())
}

func LogMigrationResults(results []*goose.MigrationResult, logger *slog.Logger) {
	for _, r := range results {
		LogMigrationResult(r, logger)
	}
}

func ApplyMigrations(ctx context.Context, provider *goose.Provider, direction string, logger *slog.Logger) error {
	switch direction {
	case "up":
		results, err := provider.Up(ctx)
		if err != nil {
			return err
		}
		LogMigrationResults(results, logger)

	case "down":
		result, err := provider.Down(ctx)
		if err != nil {
			return err
		}
		LogMigrationResult(result, logger)
	}

	return nil
}
