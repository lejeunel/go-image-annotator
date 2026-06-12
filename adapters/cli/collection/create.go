package collection

import (
	"context"
	"fmt"

	cli "github.com/lejeunel/go-image-annotator/adapters/cli"
	i "github.com/lejeunel/go-image-annotator/adapters/sqlite/infra"
	db "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos"
	afs "github.com/lejeunel/go-image-annotator/app/file-store"
	"github.com/lejeunel/go-image-annotator/config"
	clc "github.com/lejeunel/go-image-annotator/use-cases/collection/create"
)

type CreatePresenter struct {
	cli.ErrorPresenter
}

func (p CreatePresenter) Success(r clc.Response) {
	fmt.Println("created collection with name", r.Name, "and description", r.Description)
}
func (p CreatePresenter) Error(err error) {
	fmt.Println(err.Error())
}

func Create(name string, group *string, description string) {
	cfg := config.Parse()
	app := i.NewSQLiteInfra(db.NewSQLiteDB(cfg.DBPath),
		afs.NewFileStore(cfg.ArtefactDir))
	itr := clc.New(app.Collection, app.Group)
	itr.Execute(context.Background(),
		clc.Request{Name: name, Group: group, Description: description}, CreatePresenter{})

}
