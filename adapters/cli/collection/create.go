package collection

import (
	"context"
	"fmt"

	cli "github.com/lejeunel/go-image-annotator/adapters/cli"
	a "github.com/lejeunel/go-image-annotator/app/sqlite"
	"github.com/lejeunel/go-image-annotator/config"
	clc "github.com/lejeunel/go-image-annotator/use-cases/collection/create"
)

type CreatePresenter struct {
	cli.ErrorPresenter
}

func (p CreatePresenter) Success(r clc.Response) {
	fmt.Println("created collection with name", r.Name, "and description", r.Description)
}

func Create(name string, group *string, description string) {
	app := a.NewSQLiteApp(config.Parse())
	app.Itrs.Collection.Create.Execute(context.Background(),
		clc.Request{Name: name, Group: group, Description: description}, CreatePresenter{})

}
