package user

import (
	"context"
	"fmt"

	cli "github.com/lejeunel/go-image-annotator/adapters/cli"
	i "github.com/lejeunel/go-image-annotator/adapters/sqlite/infra"
	db "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos"
	afs "github.com/lejeunel/go-image-annotator/app/file-store"
	tok "github.com/lejeunel/go-image-annotator/app/token"
	"github.com/lejeunel/go-image-annotator/config"
	uc "github.com/lejeunel/go-image-annotator/use-cases/user/create"
)

type Presenter struct {
	cli.ErrorPresenter
}

func (p Presenter) Success(r uc.Response) {
	fmt.Println("created user with id", r.Id, "and admin rights", r.IsAdmin)
}
func (p Presenter) Error(err error) {
	fmt.Println(err.Error())
}

func Create(id string, isAdmin bool) {
	cfg := config.Parse()
	app := i.NewSQLiteInfra(db.NewSQLiteDB(cfg.DBPath),
		afs.NewFileStore(cfg.ArtefactDir))
	itr := uc.New(app.User, tok.NewTokenGenerator(32))
	itr.Execute(context.Background(),
		uc.Request{Id: id, IsAdmin: isAdmin}, Presenter{})
}
