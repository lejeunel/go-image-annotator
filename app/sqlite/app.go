package sqlite

import (
	"crypto/sha256"

	db "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	"github.com/lejeunel/go-image-annotator/app"
	"github.com/lejeunel/go-image-annotator/config"
	a "github.com/lejeunel/go-image-annotator/modules/annotator"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	fs "github.com/lejeunel/go-image-annotator/modules/file-store"
	ig "github.com/lejeunel/go-image-annotator/modules/ingester"
	pv "github.com/lejeunel/go-image-annotator/modules/password-validator"
	rea "github.com/lejeunel/go-image-annotator/modules/reader"
	"github.com/lejeunel/go-image-annotator/modules/scroller"
	tk "github.com/lejeunel/go-image-annotator/modules/token"
	sm "github.com/lejeunel/go-image-annotator/shared/session"
)

func NewSQLiteApp(cfg config.Config, auth auth.Authorizer) app.App {
	apiTokenGen := tk.New(cfg.ApiTokenLength)
	passwordGen := tk.New(cfg.RandomPasswordLength)
	forgottenPasswordGen := tk.New(cfg.RandomPasswordLength)
	passwordValidator := pv.New(cfg.PasswordMinEntropy)
	sqldb := db.NewSQLiteDB(cfg.SQLiteDBPath)
	fileStore := fs.NewFileStore(cfg.ArtefactDir)
	repos := NewSQLiteRepos(sqldb, fileStore)
	sessionManager := sm.NewSQLiteSessionManager(sqldb.DB, repos.User, apiTokenGen)
	scr := scroller.New(repos.Scroller)
	ingester := ig.New(repos.Image, repos.Collection, repos.Label, repos.Annotation, fileStore, sha256.New(), rea.ImageSpecsDetector{})
	itrs := NewSQLiteInteractors(
		repos,
		cfg.DefaultPageSize,
		cfg.AllowedImageFormats,
		apiTokenGen,
		passwordGen,
		forgottenPasswordGen,
		cfg.ForgotPasswordTokenExpirationMinutes,
		passwordValidator,
		apiTokenGen,
		ingester,
		auth)
	annotator := a.NewAnnotator(scr, itrs.Image.Find,
		itrs.Annotation.AddBox, itrs.Annotation.UpdateBox,
		itrs.Annotation.AddPolygon, itrs.Annotation.UpdatePolygon,
		itrs.Annotation.Delete,
		itrs.Label.FetchAll, itrs.Annotation.UpdateLabel,
		itrs.Annotation.AddImageLabel)

	return app.App{
		Itrs:           itrs,
		SessionManager: sessionManager,
		Annotator:      annotator,
	}

}
