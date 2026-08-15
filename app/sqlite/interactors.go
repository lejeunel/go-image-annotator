package sqlite

import (
	"crypto/sha256"
	itr "github.com/lejeunel/go-image-annotator/app/interactors"
	cfg "github.com/lejeunel/go-image-annotator/config"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	"log/slog"

	tra "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/transactors"
	aig "github.com/lejeunel/go-image-annotator/modules/archive-ingester"
	el "github.com/lejeunel/go-image-annotator/modules/event-logger"
	iig "github.com/lejeunel/go-image-annotator/modules/image-ingester"
	ims "github.com/lejeunel/go-image-annotator/modules/image-store"
	pv "github.com/lejeunel/go-image-annotator/modules/password-validator"
	rea "github.com/lejeunel/go-image-annotator/modules/reader"
	tk "github.com/lejeunel/go-image-annotator/modules/token"
)

func BuildInteractors(infra Infra, auth auth.Interface, logger slog.Logger, cfg cfg.Config, ts tk.TokenService) itr.Interactors {
	passwordTokenizer := tk.New(cfg.RandomPasswordLength)
	forgottenPasswordGen := tk.New(cfg.RandomPasswordLength)
	passwordValidator := pv.New(cfg.PasswordMinEntropy)

	imstore := ims.New(
		ims.Repos{
			ImageRepo:      infra.ImageRepo,
			CollectionRepo: infra.CollectionRepo,
			AnnotationRepo: infra.AnnotationRepo,
			MetaRepo:       infra.MetaRepo,
		},
		tra.NewStoreTransactor(infra.DB, infra.FilterStrParser, infra.OrderStrParser),
		infra.ImageFileStore)
	eventlogger := el.New(infra.EventRepo, el.WithMaxNumTasksPerUser(cfg.MaxNumTasksPerUser))

	imageIngester := iig.New(infra.ImageRepo, infra.CollectionRepo, infra.LabelRepo, infra.AnnotationRepo,
		tra.NewIngestionTransactor(infra.DB),
		infra.ImageFileStore, sha256.New(), rea.NewImageSpecsDetector(cfg.AllowedImageMIMETypes))
	archiveIngester := aig.New(imstore, imageIngester)

	return itr.Interactors{
		Label: NewLabelInteractors(infra.LabelRepo, cfg.DefaultPageSize, auth),
		Collection: NewCollectionInteractors(
			infra.DB,
			infra.CollectionRepo,
			infra.ImageRepo,
			infra.AnnotationRepo,
			infra.GroupRepo,
			imstore,
			eventlogger,
			logger,
			cfg.DefaultPageSize,
			auth,
		),
		Image: NewImageInteractors(
			infra.ImageRepo,
			infra.CollectionRepo,
			infra.AnnotationRepo,
			imstore,
			infra.ImageFileStore,
			infra.TempFileStore,
			imageIngester,
			archiveIngester,
			infra.FilterStrParser,
			infra.OrderingStrConverter,
			int64(cfg.MaxArchiveMB),
			eventlogger,
			logger,
			cfg.DefaultPageSize,
			cfg.MaxPageSize,
			auth,
		),
		User: NewUserInteractors(infra.UserRepo, infra.GroupRepo, infra.RoleRepo,
			ts,
			forgottenPasswordGen,
			passwordValidator,
			passwordTokenizer,
			passwordTokenizer,
			cfg.ForgotPasswordTokenExpirationMinutes,
			forgottenPasswordGen, auth),
		Annotation: NewAnnotationInteractors(imstore, infra.ImageRepo, infra.LabelRepo, infra.AnnotationRepo, auth),
		Group:      NewGroupInteractors(infra.GroupRepo, auth),
		Role:       NewRoleInteractors(infra.RoleRepo, auth),
		Bootstrap: NewBootstrapInteractor(
			infra.UserRepo,
			infra.RoleRepo,
			infra.PolicyFileStore,
			passwordTokenizer,
			passwordValidator,
		),
		Policy:   NewPolicyInteractors(infra.PolicyFileStore, auth),
		Metadata: NewMetadataInteractors(infra.MetaRepo, infra.CollectionRepo, infra.ImageRepo, auth),
		Log:      NewLogInteractors(eventlogger),
	}

}
