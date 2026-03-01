package app

import (
	"context"
	"database/sql"
	au "datahub/app/authorizer"
	in "datahub/app/ingester"
	m "datahub/app/migrations"
	c "datahub/config"
	pro "datahub/domain/annotation_profiles"
	clc "datahub/domain/collections"
	im "datahub/domain/images"
	lbl "datahub/domain/labels"
	loc "datahub/domain/locations"
	clk "github.com/jonboulle/clockwork"
	goose "github.com/pressly/goose/v3"
	"log/slog"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

type NoopLoggingHandler struct{}

func (h *NoopLoggingHandler) Enabled(context.Context, slog.Level) bool  { return false }
func (h *NoopLoggingHandler) Handle(context.Context, slog.Record) error { return nil }
func (h *NoopLoggingHandler) WithAttrs([]slog.Attr) slog.Handler        { return h }
func (h *NoopLoggingHandler) WithGroup(string) slog.Handler             { return h }

type App struct {
	Clock       clk.Clock
	Images      *im.Service
	Collections *clc.Service
	Locations   *loc.Service
	Profiles    *pro.Service
	Labels      *lbl.Service
	Ingestion   *in.Service
	Config      *c.Config
	Authorizer  *au.Authorizer
	Logger      *slog.Logger
}

func NewApp(cfg *c.Config, clock clk.Clock, verbose int) (App, *sql.DB, *goose.Provider, context.Context) {

	logger := buildLogger(cfg, verbose)

	db, migrationProvider, err := buildDbAndMigrations(cfg, logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	repos := buildRepos(db, logger)
	auth := buildAuthorizer(cfg, logger)
	kvStore := buildKVStore(cfg, logger)

	locationService := loc.NewLocationService(repos.LocationRepo, cfg.MaxPageSize, cfg.MaxPageSize,
		auth, logger, clock)
	labelsService := lbl.NewLabelService(repos.LabelRepo, cfg.MaxPageSize, cfg.MaxPageSize, auth, logger, clock)
	profileService := pro.NewAnnotationProfileService(repos.AnnotationProfileRepo, labelsService, logger, auth, clock)
	collectionService := clc.NewCollectionService(repos.CollectionRepo, profileService, labelsService, cfg.MaxPageSize, cfg.MaxPageSize,
		auth, logger, clock)
	imageAnnotationService := im.NewAnnotationService(repos.AnnotationRepo,
		collectionService,
		labelsService, logger, auth, clock)
	imageService := im.NewImageService(&kvStore, repos.ImageRepo, imageAnnotationService,
		locationService,
		collectionService, cfg.MaxPageSize, logger,
		cfg.AllowedImageTypes, auth, clock)
	ingestionService := in.NewIngestionService(imageService, collectionService, labelsService,
		locationService, logger, auth)

	ctx := context.Background()

	app := App{
		Clock:       clock,
		Images:      imageService,
		Collections: collectionService,
		Locations:   locationService,
		Profiles:    profileService,
		Labels:      labelsService,
		Ingestion:   ingestionService,
		Config:      cfg,
		Authorizer:  auth,
		Logger:      logger}

	return app, db.DB, migrationProvider, ctx

}

func NewTestApp(t *testing.T, isAdmin bool) (App, *clk.FakeClock, context.Context) {
	cfg := &c.Config{MaxPageSize: 999, AllowedImageTypes: []string{"thermal", "rgb", "gray"},
		Mode: "test"}

	clock := clk.NewFakeClock()
	app, _, migrationProvider, ctx := NewApp(cfg, clock, 0)
	if err := m.ApplyMigrations(ctx, migrationProvider, "up", app.Logger); err != nil {
		panic(err)
	}

	if isAdmin == true {
		ctx = context.WithValue(ctx, "entitlements", "admin")
		ctx = context.WithValue(ctx, "groups", "mygroup")
	}
	ctx = context.WithValue(ctx, "username", "test-user")
	ctx = context.WithValue(ctx, "email", "test-user@mail.com")

	return app, clock, ctx
}
