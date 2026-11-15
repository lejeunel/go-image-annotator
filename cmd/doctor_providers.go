package cmd

import (
	"log/slog"

	db "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	app "github.com/lejeunel/go-image-annotator/app/sqlite"
	"github.com/lejeunel/go-image-annotator/config"
	fs "github.com/lejeunel/go-image-annotator/modules/file-store"
	ntf "github.com/lejeunel/go-image-annotator/modules/notifier"
	pv "github.com/lejeunel/go-image-annotator/modules/password-validator"
	dr "github.com/lejeunel/go-image-annotator/use-cases/doctor"
)

type FromEnvConfigProvider struct{}

func (p FromEnvConfigProvider) Get() (*config.Config, error) {
	return config.Parse()
}

type PasswordValidatorProvider struct{}

func (p PasswordValidatorProvider) Get(minEntropy int) pv.PasswordValidator {
	return pv.New(minEntropy)
}

type FileStoreProvider struct{}

func (p FileStoreProvider) Get(c config.Config) (fs.FileStore, error) {
	return fs.Build(c)
}

type DBStoreProvider struct{}

func (p DBStoreProvider) Get(path string) (fs.FileStore, error) {
	return fs.NewLocalFileStore(path)
}

type ImageStoreProvider struct{}

func (p ImageStoreProvider) Get(c config.Config) (dr.ImageStore, error) {
	fileStore, err := fs.Build(c)
	if err != nil {
		return nil, err
	}

	imageStore := app.BuildImageStoreFromInfra(app.BuildInfra(c.DBPath, fileStore))
	return &imageStore, nil
}

func (p ImageStoreProvider) IsInstantiated(c config.Config) bool {
	if db.DBExists(c.DBPath) {
		return true
	} else {
		return false
	}
}

type NotifierProvider struct{}

func (p NotifierProvider) Get(c config.Config) (ntf.Notifier, error) {
	l := slog.New(slog.DiscardHandler)
	return ntf.NewNotifierWithVoidDefault(c.SMTPHost, c.SMTPPort, c.SMTPUsername, c.SMTPPassword,
		*l)
}
