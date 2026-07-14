package sqlite

import (
	db "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	"github.com/lejeunel/go-image-annotator/app"
	"github.com/lejeunel/go-image-annotator/config"
	a "github.com/lejeunel/go-image-annotator/modules/annotator"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	fs "github.com/lejeunel/go-image-annotator/modules/file-store"
	pv "github.com/lejeunel/go-image-annotator/modules/password-validator"
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
	repos := NewSQLiteRepos(sqldb,
		fs.NewFileStore(cfg.ArtefactDir))
	sessionManager := sm.NewSQLiteSessionManager(sqldb.DB, repos.User, apiTokenGen)
	scr := scroller.New(repos.Scroller)
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
