package app

import (
	"context"
	"log/slog"
	"os"

	itrs "github.com/lejeunel/go-image-annotator/app/interactors"
	"github.com/lejeunel/go-image-annotator/config"
	r "github.com/lejeunel/go-image-annotator/entities/role"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	a "github.com/lejeunel/go-image-annotator/modules/annotator"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	q "github.com/lejeunel/go-image-annotator/modules/query"
	s "github.com/lejeunel/go-image-annotator/shared/session"
	bst "github.com/lejeunel/go-image-annotator/use-cases/bootstrap"
)

type ImageFilterDocumenter interface {
	DescribeFields() []q.FieldDescription
	Examples() []string
}
type ImageSortDocumenter interface {
	DescribeOrderingFields() []q.FieldDescription
	Examples() []string
}

type DBManager interface {
	Status(context.Context)
	Up(context.Context)
	Down(context.Context)
	Init(context.Context)
}

type App struct {
	config.Config
	Itrs itrs.Interactors
	s.SessionManager
	a.Annotator
	ImageFilterDocumenter
	ImageSortDocumenter
	auth.Authorizer
	DBManager
	slog.Logger
}

func NewApp(
	cfg config.Config,
	itrs itrs.Interactors,
	sm s.SessionManager,
	an a.Annotator,
	fd ImageFilterDocumenter,
	sd ImageSortDocumenter,
	auth auth.Authorizer,
	dbm DBManager,
	logger slog.Logger,
) App {
	return App{
		Config:                cfg,
		Itrs:                  itrs,
		SessionManager:        sm,
		Annotator:             an,
		ImageFilterDocumenter: fd,
		ImageSortDocumenter:   sd,
		Authorizer:            auth,
		DBManager:             dbm,
		Logger:                logger,
	}
}

type InitialAdminPresenter struct {
	slog.Logger
}

func (p InitialAdminPresenter) SuccessBootstrap(r bst.Response) {
	if !r.Skipped {
		p.Logger.Info("successfully bootstrapped application with initial admin")
	}
}

func (p InitialAdminPresenter) Error(err error) {
	p.Logger.Error("failed bootstrapping application", "error", err)
	os.Exit(1)
}

func BootstrapInitialAdmin(itr bst.Interactor, email, password string, logger slog.Logger) {
	role := r.NewRole(r.NewRoleId(), "admin", r.WithMethods([]string{"*"}))
	user := u.NewUser("anonymous", u.WithRoles([]r.Role{role}))
	ctx := u.AppendUserToContext(context.Background(), user)
	pres := InitialAdminPresenter{logger}
	itr.Execute(ctx, bst.Request{InitialAdminEmail: email, InitialAdminPassword: password}, pres)
}
