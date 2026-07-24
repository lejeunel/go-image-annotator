package app

import (
	"context"
	"errors"
	"log/slog"

	u "github.com/lejeunel/go-image-annotator/entities/user"
	a "github.com/lejeunel/go-image-annotator/modules/annotator"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	s "github.com/lejeunel/go-image-annotator/shared/session"
	createusr "github.com/lejeunel/go-image-annotator/use-cases/user/create"
)

type App struct {
	Itrs           Interactors
	SessionManager s.MySessionManager
	a.Annotator
}

type InitialAdminPresenter struct {
	slog.Logger
}

func (p InitialAdminPresenter) Success(r createusr.Response) {
	p.Logger.Info("created initial admin user", "id", r.Id)
}
func (p InitialAdminPresenter) Error(err error) {
	if errors.Is(err, e.ErrDuplicate) {
		return
	}
	p.Logger.Error("creating initial admin user", "error", err)
}

func MaybeCreateInitialAdmin(itr createusr.Interactor, email, password string) {
	initAdminUser := u.NewUser("admin", u.WithAdmin(true))
	itr.Execute(
		u.AppendUserToContext(context.Background(), initAdminUser),
		createusr.Request{Id: email, Password: &password, IsAdmin: true},
		InitialAdminPresenter{})
}
