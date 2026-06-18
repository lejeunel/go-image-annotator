package sqlite

import (
	db "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	"github.com/lejeunel/go-image-annotator/app"
	"github.com/lejeunel/go-image-annotator/config"
	a "github.com/lejeunel/go-image-annotator/modules/annotator"
	"github.com/lejeunel/go-image-annotator/modules/auth"
	fs "github.com/lejeunel/go-image-annotator/modules/file-store"
	"github.com/lejeunel/go-image-annotator/modules/scroller"
	tok "github.com/lejeunel/go-image-annotator/modules/token"
	ip "github.com/lejeunel/go-image-annotator/shared/identity_provider"
	sm "github.com/lejeunel/go-image-annotator/shared/session"
)

func NewSQLiteApp(cfg config.Config, auth auth.Auth) app.App {
	tg := tok.NewTokenGenerator(cfg.TokenLength)
	sqldb := db.NewSQLiteDB(cfg.SQLiteDBPath)
	repos := NewSQLiteRepos(sqldb,
		fs.NewFileStore(cfg.ArtefactDir))
	sessionManager := sm.NewSQLiteSessionManager(sqldb.DB, repos.User, tg)
	identityProvider := ip.NewGothIdentityHandler(sessionManager)
	scr := scroller.New(repos.Scroller)
	itrs := NewSQLiteInteractors(repos,
		cfg.DefaultPageSize, cfg.AllowedImageFormats, tg, auth)
	annotator := a.NewAnnotator(scr, itrs.Image.Read,
		itrs.Annotation.AddBox, itrs.Annotation.UpdateBox,
		itrs.Annotation.AddPolygon, itrs.Annotation.UpdatePolygon,
		itrs.Annotation.Delete,
		itrs.Label.FetchAll, itrs.Annotation.UpdateLabel,
		itrs.Annotation.AddImageLabel)

	return app.App{
		Itrs:           itrs,
		SessionManager: sessionManager,
		OAuthHandler:   identityProvider,
		Annotator:      annotator,
	}

}
