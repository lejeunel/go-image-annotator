package user

import (
	"context"
	"fmt"

	cli "github.com/lejeunel/go-image-annotator/adapters/cli"
	a "github.com/lejeunel/go-image-annotator/app/sqlite"
	"github.com/lejeunel/go-image-annotator/config"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	uc "github.com/lejeunel/go-image-annotator/use-cases/user/create"
)

type Presenter struct {
	cli.ErrorPresenter
}

func (p Presenter) Success(r uc.Response) {
	fmt.Println("created user with id", r.Id, "and admin rights", r.IsAdmin)
}

func Create(id string, isAdmin bool) {
	app := a.NewSQLiteApp(config.Parse(), auth.NewVoidAuth())
	app.Itrs.User.Create.Execute(context.Background(),
		uc.Request{Id: id, IsAdmin: isAdmin}, Presenter{})
}
