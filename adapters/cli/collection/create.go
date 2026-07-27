package collection

import (
	"context"

	cli "github.com/lejeunel/go-image-annotator/adapters/cli"
	a "github.com/lejeunel/go-image-annotator/app/sqlite"
	"github.com/lejeunel/go-image-annotator/config"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	clc "github.com/lejeunel/go-image-annotator/use-cases/collection/create"
)

type CreatePresenter struct {
	cli.ErrorPresenter
}

func (p CreatePresenter) Success(r clc.Response) {
	p.Logger.Info("created collection", "name", r.Name, "description", r.Description)
}

func Create(ctx context.Context, name string, group *string, description string) {
	auth := auth.NewDefault()
	app := a.NewSQLiteApp(config.Parse(), &auth)
	app.Itrs.Collection.Create.Execute(ctx,
		clc.Request{Name: name, Group: group, Description: description},
		CreatePresenter{cli.NewErrorPresenter()})

}
