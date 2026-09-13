package sqlite

import (
	"log/slog"

	"github.com/lejeunel/go-image-annotator/app"
	"github.com/lejeunel/go-image-annotator/config"
	a "github.com/lejeunel/go-image-annotator/modules/annotator"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	fs "github.com/lejeunel/go-image-annotator/modules/file-store"
	tk "github.com/lejeunel/go-image-annotator/modules/token"
)

func NewApp(cfg config.Config, auth auth.Interface, logger slog.Logger) app.App {
	imageStore, err := fs.Build(cfg, logger)
	if err != nil {
		panic(err)
	}

	infra := BuildInfra(cfg.LocalArtefactPath, imageStore)
	apiTokenGen := tk.New(cfg.ApiTokenLength)
	itrs := BuildInteractors(infra, auth, logger, cfg, apiTokenGen)
	sessionManager := NewSessionManager(infra.DB.DB, infra.UserRepo, apiTokenGen)

	annotator := a.NewAnnotator(itrs.Image.Scroll, itrs.Image.Find,
		itrs.Annotation.AddBox, itrs.Annotation.UpdateBox,
		itrs.Annotation.AddPolygon, itrs.Annotation.UpdatePolygon,
		itrs.Annotation.Delete,
		itrs.Label.FetchAll, itrs.Annotation.UpdateLabel,
		itrs.Annotation.AddImageLabel, itrs.Metadata.Add, itrs.Metadata.List,
		itrs.Metadata.Read, itrs.Metadata.Delete,
	)

	return app.NewApp(itrs, sessionManager, annotator, infra.FilterParser, infra.OrderParser)
}
