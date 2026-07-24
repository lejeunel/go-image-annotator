package app

import (
	"context"
	"errors"
	"log/slog"
	"os"

	u "github.com/lejeunel/go-image-annotator/entities/user"
	a "github.com/lejeunel/go-image-annotator/modules/annotator"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	s "github.com/lejeunel/go-image-annotator/shared/session"
	createrl "github.com/lejeunel/go-image-annotator/use-cases/role/create"
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

func (p InitialAdminPresenter) SuccessCreateUser(r createusr.Response) {
	p.Logger.Info("created initial admin user", "id", r.Id)
}
func (p InitialAdminPresenter) SuccessCreateRole(r createrl.Response) {
	p.Logger.Info("created admin role")
}
func (p InitialAdminPresenter) Error(err error) {
	if errors.Is(err, e.ErrDuplicate) {
		return
	}
	p.Logger.Error("creating initial admin user", "error", err)
	os.Exit(1)
}

func BootstrapInitialAdmin(userCreator createusr.Interactor, roleCreator createrl.Interactor, email, password string, logger slog.Logger) {
	roles := []string{"admin"}
	user := u.NewUser(email, u.WithRoles(roles))
	ctx := u.AppendUserToContext(context.Background(), user)
	pres := InitialAdminPresenter{logger}
	roleCreator.Execute(ctx, createrl.Request{Name: "admin", Description: ""}, pres)
	userCreator.Execute(ctx,
		createusr.Request{Id: email, Password: &password, Roles: roles}, pres)
}
