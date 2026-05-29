package collection

import (
	"fmt"
	"github.com/lejeunel/go-image-annotator/config"

	cli "github.com/lejeunel/go-image-annotator/adapters/cli"
	"github.com/lejeunel/go-image-annotator/infra"
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

func Create(name, description string) {
	cfg := config.Parse()
	app := infra.NewSQLiteInfra(cfg.DBPath, cfg.ArtefactDir)
	itr := clc.NewInteractor(app.CollectionRepo)
	itr.Execute(clc.Request{Name: name, Description: description}, CreatePresenter{})

}
