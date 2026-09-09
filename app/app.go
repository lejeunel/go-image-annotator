package app

import (
	"context"
	"log/slog"
	"os"

	itrs "github.com/lejeunel/go-image-annotator/app/interactors"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	a "github.com/lejeunel/go-image-annotator/modules/annotator"
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
type App struct {
	Itrs itrs.Interactors
	s.SessionManager
	a.Annotator
	ImageFilterDocumenter
	ImageSortDocumenter
}

func NewApp(itrs itrs.Interactors, sm s.SessionManager, an a.Annotator, fd ImageFilterDocumenter, sd ImageSortDocumenter) App {
	return App{
		Itrs:                  itrs,
		SessionManager:        sm,
		Annotator:             an,
		ImageFilterDocumenter: fd,
		ImageSortDocumenter:   sd,
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
	user := u.NewUser("anonymous", u.WithRoles([]string{"admin"}))
	ctx := u.AppendUserToContext(context.Background(), user)
	pres := InitialAdminPresenter{logger}
	itr.Execute(ctx, bst.Request{InitialAdminEmail: email, InitialAdminPassword: password}, pres)
}
