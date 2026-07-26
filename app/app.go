package app

import (
	"context"
	"log/slog"
	"os"

	itrs "github.com/lejeunel/go-image-annotator/app/interactors"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	a "github.com/lejeunel/go-image-annotator/modules/annotator"
	s "github.com/lejeunel/go-image-annotator/shared/session"
	bst "github.com/lejeunel/go-image-annotator/use-cases/bootstrap"
)

type App struct {
	Itrs           itrs.Interactors
	SessionManager s.MySessionManager
	a.Annotator
}

type InitialAdminPresenter struct {
	slog.Logger
}

func (p InitialAdminPresenter) SuccessBootstrap() {
	p.Logger.Info("successfully bootstrapped application with initial admin")
}
func (p InitialAdminPresenter) Error(err error) {
	p.Logger.Error("failed bootstrapping application", "error", err)
	os.Exit(1)
}

func BootstrapInitialAdmin(itr bst.Interactor, email, password string, logger slog.Logger) {
	user := u.NewUser("anonymous", u.WithRoles([]string{"admin"}))
	ctx := u.AppendUserToContext(context.Background(), user)
	pres := InitialAdminPresenter{logger}
	itr.Execute(ctx, bst.Request{InitialAdminEmail: email, InitialAdminPassword: password}, pres)
}
