package sqlite

import (
	"log/slog"

	sqldb "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	"github.com/lejeunel/go-image-annotator/app"
	"github.com/lejeunel/go-image-annotator/config"

	a "github.com/lejeunel/go-image-annotator/modules/annotator"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	fs "github.com/lejeunel/go-image-annotator/modules/file-store"
	tk "github.com/lejeunel/go-image-annotator/modules/token"
)

func NewAppFromEnv(logger *slog.Logger) app.App {
	cfg := config.Parse()
	return NewApp(cfg, *logger)
}

func NewApp(cfg config.Config, logger slog.Logger) app.App {
	imageStore, err := fs.Build(cfg, logger)
	if err != nil {
		panic(err)
	}

	infra := BuildInfra(cfg.LocalArtefactPath, imageStore)

	auth := auth.New(auth.ValidMethods, infra.RoleRepo)
	apiTokenGen := tk.New(cfg.ApiTokenLength)
	itrs := BuildInteractors(infra, auth, logger, cfg, apiTokenGen)
	sessionManager := NewSessionManager(infra.DB.DB, infra.UserStore, apiTokenGen)

	annotator := a.NewAnnotator(itrs.Image.Scroll, itrs.Image.Find,
		itrs.Annotation.AddBox, itrs.Annotation.UpdateBox,
		itrs.Annotation.AddPolygon, itrs.Annotation.UpdatePolygon,
		itrs.Annotation.Delete,
		itrs.Annotation.PickLabel, itrs.Annotation.UpdateLabel,
		itrs.Annotation.AddImageLabel, itrs.Metadata.Add, itrs.Metadata.List,
		itrs.Metadata.Read, itrs.Metadata.Delete,
	)

	dbManager := sqldb.NewSQLiteDBManager(infra.DB.DB, &logger)
	return app.NewApp(
		cfg,
		itrs,
		sessionManager,
		annotator,
		infra.FilterParser,
		infra.OrderParser,
		auth,
		&dbManager,
		logger,
	)
}
