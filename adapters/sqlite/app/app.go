package app

import (
	infra "github.com/lejeunel/go-image-annotator/adapters/sqlite/infra"
	i "github.com/lejeunel/go-image-annotator/adapters/sqlite/interactors"
	db "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos"
	fs "github.com/lejeunel/go-image-annotator/app/file-store"
	tok "github.com/lejeunel/go-image-annotator/app/token"
	"github.com/lejeunel/go-image-annotator/config"
	u "github.com/lejeunel/go-image-annotator/use-cases"
)

type SQLiteApp struct {
	Infra infra.SQLiteInfra
	Itrs  u.Interactors
	tok.TokenGenerator
}

func NewSQLiteApp(cfg config.Config) SQLiteApp {
	infra := infra.NewSQLiteInfra(db.NewSQLiteDB(cfg.DBPath),
		fs.NewFileStore(cfg.ArtefactDir))
	tg := tok.NewTokenGenerator(cfg.TokenLength)
	return SQLiteApp{
		Infra:          infra,
		Itrs:           i.NewSQLiteInteractors(infra, cfg.DefaultPageSize, cfg.AllowedImageFormats, tg),
		TokenGenerator: tg,
	}

}
