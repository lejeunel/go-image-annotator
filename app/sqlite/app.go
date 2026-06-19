package sqlite

import (
	db "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	"github.com/lejeunel/go-image-annotator/app"
	"github.com/lejeunel/go-image-annotator/config"
	a "github.com/lejeunel/go-image-annotator/modules/annotator"
	"github.com/lejeunel/go-image-annotator/modules/auth"
	au "github.com/lejeunel/go-image-annotator/modules/authentifier"
	fs "github.com/lejeunel/go-image-annotator/modules/file-store"
	"github.com/lejeunel/go-image-annotator/modules/scroller"
	ip "github.com/lejeunel/go-image-annotator/shared/identity_provider"
	sm "github.com/lejeunel/go-image-annotator/shared/session"
)

func NewSQLiteApp(cfg config.Config, auth auth.Auth) app.App {
	tg := au.New(cfg.TokenLength)
	pg := au.New(cfg.RandomPasswordLength)
	sqldb := db.NewSQLiteDB(cfg.SQLiteDBPath)
	repos := NewSQLiteRepos(sqldb,
		fs.NewFileStore(cfg.ArtefactDir))
	sessionManager := sm.NewSQLiteSessionManager(sqldb.DB, repos.User, tg)
	identityProvider := ip.NewGothIdentityHandler(sessionManager)
	scr := scroller.New(repos.Scroller)
	itrs := NewSQLiteInteractors(repos,
		cfg.DefaultPageSize, cfg.AllowedImageFormats, tg, pg, auth)
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
