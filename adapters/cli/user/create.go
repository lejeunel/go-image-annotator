package user

import (
	"context"
	"fmt"

	cli "github.com/lejeunel/go-image-annotator/adapters/cli"
	"github.com/lejeunel/go-image-annotator/adapters/sqlite"
	tg "github.com/lejeunel/go-image-annotator/app/token-generator"
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
	app := infra.NewSQLiteInfra(cfg.DBPath, cfg.ArtefactDir)
	itr := uc.New(app.User, tg.NewTokenGenerator(32))
	itr.Execute(context.Background(),
		uc.Request{Id: id, IsAdmin: isAdmin}, Presenter{})
}
