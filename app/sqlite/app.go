package sqlite

import (
	"crypto/sha256"
	"fmt"
	"log/slog"

	db "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	an "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/annotation"
	clc "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/collection"
	ev "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/event"
	grp "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/group"
	im "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/image"
	lbl "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/label"
	md "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/metadata"
	r "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/role"
	scr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/scroll"
	tra "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/transactors"
	usr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/user"
	"github.com/lejeunel/go-image-annotator/app"
	itr "github.com/lejeunel/go-image-annotator/app/interactors"
	"github.com/lejeunel/go-image-annotator/config"
	a "github.com/lejeunel/go-image-annotator/modules/annotator"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	el "github.com/lejeunel/go-image-annotator/modules/event-logger"
	fs "github.com/lejeunel/go-image-annotator/modules/file-store"
	im_store "github.com/lejeunel/go-image-annotator/modules/image-store"
	ig "github.com/lejeunel/go-image-annotator/modules/ingester"
	pv "github.com/lejeunel/go-image-annotator/modules/password-validator"
	rea "github.com/lejeunel/go-image-annotator/modules/reader"
	"github.com/lejeunel/go-image-annotator/modules/scroller"
	tk "github.com/lejeunel/go-image-annotator/modules/token"
	sm "github.com/lejeunel/go-image-annotator/shared/session"
)

func NewSQLiteApp(cfg config.Config, auth auth.Interface, logger slog.Logger) app.App {
	apiTokenGen := tk.New(cfg.ApiTokenLength)
	passwordTokenizer := tk.New(cfg.RandomPasswordLength)
	forgottenPasswordGen := tk.New(cfg.RandomPasswordLength)
	passwordValidator := pv.New(cfg.PasswordMinEntropy)
	db := db.NewSQLiteDB(fmt.Sprintf("%v/%v", cfg.ArtefactPath, "db.sqlite"))
	imrepo := im.NewSQLiteImageRepo(db)
	anrepo := an.NewSQLiteAnnotationRepo(db)
	clrepo := clc.NewSQLiteCollectionRepo(db)
	lbrepo := lbl.NewSQLiteLabelRepo(db)
	grprepo := grp.NewSQLiteGroupRepo(db)
	rlrepo := r.NewSQLiteRoleRepo(db)
	usrrepo := usr.NewSQLiteUserRepo(db)
	eventrepo := ev.NewSQLiteEventRepo(db)
	metadatarepo := md.NewSQLiteMetaRepo(db)
	imageFileStore := fs.NewLocalFileStore(fmt.Sprintf("%v/%v", cfg.ArtefactPath, "images"))
	policyFileStore := fs.NewLocalFileStore(fmt.Sprintf("%v/%v", cfg.ArtefactPath, "assets"))
	imstore := im_store.New(imrepo, clrepo, anrepo, metadatarepo, imageFileStore)
	scrrepo := scr.NewSQLiteScrollerRepo(db)
	eventlogger := el.New(eventrepo, el.WithMaxNumTasksPerUser(cfg.MaxNumTasksPerUser))

	sessionManager := sm.NewSQLiteSessionManager(db.DB, usrrepo, apiTokenGen)
	scr := scroller.New(scrrepo)
	ingester := ig.New(imrepo, clrepo, lbrepo, anrepo,
		tra.NewIngestionTransactor(db),
		imageFileStore, sha256.New(), rea.NewImageSpecsDetector(cfg.AllowedImageMIMETypes))
	itrs := itr.Interactors{
		Label: NewSQLiteLabelInteractors(lbrepo, cfg.DefaultPageSize, auth),
		Collection: NewSQLiteCollectionInteractors(
			db,
			clrepo,
			imrepo,
			anrepo,
			grprepo,
			imstore,
			eventlogger,
			logger,
			cfg.DefaultPageSize,
			auth,
		),
		Image: NewSQLiteImageInteractors(
			imrepo,
			clrepo,
			anrepo,
			imstore,
			imageFileStore,
			ingester,
			cfg.DefaultPageSize,
			auth,
		),
		User: NewSQLiteUserInteractors(usrrepo, grprepo, rlrepo,
			apiTokenGen,
			forgottenPasswordGen,
			passwordValidator,
			passwordTokenizer,
			passwordTokenizer,
			cfg.ForgotPasswordTokenExpirationMinutes,
			forgottenPasswordGen, auth),
		Annotation: NewSQLiteAnnotationInteractors(imstore, imrepo, lbrepo, anrepo, auth),
		Group:      NewSQLiteGroupInteractors(grprepo, auth),
		Role:       NewSQLiteRoleInteractors(rlrepo, auth),
		Bootstrap: NewSQLiteBootstrapInteractor(
			usrrepo,
			rlrepo,
			policyFileStore,
			passwordTokenizer,
			passwordValidator,
		),
		Policy:   NewSQLitePolicyInteractors(policyFileStore, auth),
		Metadata: NewSQLiteMetadataInteractors(metadatarepo, clrepo, imrepo, auth),
		Log:      NewSQLiteLogInteractors(eventlogger),
	}
	annotator := a.NewAnnotator(scr, itrs.Image.Find,
		itrs.Annotation.AddBox, itrs.Annotation.UpdateBox,
		itrs.Annotation.AddPolygon, itrs.Annotation.UpdatePolygon,
		itrs.Annotation.Delete,
		itrs.Label.FetchAll, itrs.Annotation.UpdateLabel,
		itrs.Annotation.AddImageLabel, itrs.Metadata.Add, itrs.Metadata.List,
		itrs.Metadata.Read, itrs.Metadata.Delete,
	)

	return app.NewApp(itrs, sessionManager, annotator)
}
