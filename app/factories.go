package app

import (
	au "datahub/app/authorizer"
	m "datahub/app/migrations"
	"datahub/assets"
	c "datahub/config"
	pro "datahub/domain/annotation_profiles"
	clc "datahub/domain/collections"
	im "datahub/domain/images"
	lbl "datahub/domain/labels"
	loc "datahub/domain/locations"
	"fmt"
	"log/slog"
	"os"

	"database/sql"
	"github.com/jmoiron/sqlx"
	goose "github.com/pressly/goose/v3"
	"io/fs"
)

// here we inject SQL migration routines written in Go
// that are specific to each dialect (sqlite and postgreSQL)
func buildMigrationProvider(db *sql.DB, mode string) (*goose.Provider, error) {
	var gooseDialect goose.Dialect
	switch mode {
	case "test", "dev":
		gooseDialect = goose.DialectSQLite3
	case "prod":
		gooseDialect = goose.DialectPostgres
	}

	var collectionProfileMigration *goose.Migration
	switch mode {
	case "test", "dev":
		collectionProfileMigration = goose.NewGoMigration(202507300832,
			&goose.GoFunc{RunTx: m.Up202507300832_sqlite},
			&goose.GoFunc{RunTx: m.Down202507300832_sqlite})
	case "prod":
		collectionProfileMigration = goose.NewGoMigration(202507300832,
			&goose.GoFunc{RunTx: m.Up202507300832_pgsql},
			&goose.GoFunc{RunTx: m.Down202507300832_pgsql})
	}

	migrationsFSsub, err := fs.Sub(assets.MigrationsFS, "migrations")
	if err != nil {
		return nil, err
	}
	provider, err := goose.NewProvider(gooseDialect, db, migrationsFSsub,
		goose.WithGoMigrations(collectionProfileMigration))
	if err != nil {
		return nil, err
	}
	return provider, nil

}

func buildDbAndMigrations(cfg *c.Config, logger *slog.Logger) (*sqlx.DB, *goose.Provider, error) {

	baseErrMsg := "building database and migration provider"
	logger.Info(fmt.Sprintf("Opening db in mode %v", cfg.Mode))
	var db *sqlx.DB

	switch cfg.Mode {
	case "test":
		db = NewSQLiteConnection(":memory:", logger)
	case "dev":
		db = NewSQLiteConnection(cfg.LocalPath+"/db.sqlite", logger)
	case "prod":
		db = NewPostgresConnection(cfg, logger)
	default:
		return nil, nil, fmt.Errorf("%v: mode: %v", baseErrMsg, cfg.Mode)
	}
	provider, err := buildMigrationProvider(db.DB, cfg.Mode)
	if err != nil {
		return nil, nil, fmt.Errorf("%v: %w", baseErrMsg, err)
	}
	return db, provider, nil

}

type Repos struct {
	ImageRepo             im.ImageRepo
	AnnotationRepo        im.AnnotationRepo
	AnnotationProfileRepo pro.AnnotationProfileRepo
	LabelRepo             lbl.Repo
	CollectionRepo        clc.CollectionRepo
	LocationRepo          loc.LocationRepo
}

func buildRepos(cfg *c.Config, db *sqlx.DB, logger *slog.Logger) Repos {
	useSQLite := cfg.Mode == "test" || cfg.Mode == "dev"

	if useSQLite {
		return Repos{
			ImageRepo:             im.NewSQLiteImageRepo(db, logger),
			AnnotationRepo:        im.NewSQLiteAnnotationRepo(db),
			AnnotationProfileRepo: pro.NewSQLiteAnnotationProfileRepo(db),
			LabelRepo:             lbl.NewSQLiteLabelRepo(db),
			CollectionRepo:        clc.NewSQLiteCollectionRepo(db),
			LocationRepo:          loc.NewSQLiteLocationRepo(db),
		}
	}

	return Repos{
		ImageRepo:             im.NewPostgreSQLImageRepo(db, logger),
		AnnotationRepo:        im.NewPostgreSQLAnnotationRepo(db),
		AnnotationProfileRepo: pro.NewPostgreSQLAnnotationProfileRepo(db),
		LabelRepo:             lbl.NewPostgreSQLLabelRepo(db),
		CollectionRepo:        clc.NewPostgreSQLCollectionRepo(db),
		LocationRepo:          loc.NewPostgreSQLLocationRepo(db),
	}
}

func buildLogger(cfg *c.Config, verbose int) *slog.Logger {
	logLevel := new(slog.LevelVar)
	switch verbose {
	case 0:
		logLevel.Set(slog.LevelWarn)
	case 1:
		logLevel.Set(slog.LevelInfo)
	case 2:
		logLevel.Set(slog.LevelDebug)
	}

	if cfg.Mode == "test" {
		return slog.New(&NoopLoggingHandler{})
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
}

func buildAuthorizer(cfg *c.Config, logger *slog.Logger) *au.Authorizer {
	if cfg.Mode == "dev" {
		testingGroups := []string{"foxstream", "extern"}
		logger.Info(fmt.Sprintf("Using dummy authorizer with entitlements: %v, and groups: %v",
			cfg.TestingEntitlements, testingGroups))
		return &au.Authorizer{
			IdentityProvider: &au.TestingIdentityProvider{
				Entitlements_: cfg.TestingEntitlements,
				Groups_:       testingGroups,
				Username_:     "test-user",
				Email_:        "test-user@mail.com",
			},
		}
	}
	return au.NewAuthorizer()
}

func buildKVStore(cfg *c.Config, logger *slog.Logger) im.KeyValueStoreClient {
	if cfg.Mode == "test" {
		return im.NewMockKVStoreClient()
	}
	path := cfg.LocalPath + "/store"
	logger.Info(fmt.Sprintf("File store at %v", path))
	return im.NewFSKVStoreClient(path)
}
