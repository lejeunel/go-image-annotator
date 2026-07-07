package app

import (
	"context"
	"errors"
	"fmt"

	u "github.com/lejeunel/go-image-annotator/entities/user"
	a "github.com/lejeunel/go-image-annotator/modules/annotator"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	i "github.com/lejeunel/go-image-annotator/shared/identity_provider"
	s "github.com/lejeunel/go-image-annotator/shared/session"
	createusr "github.com/lejeunel/go-image-annotator/use-cases/user/create"
)

type App struct {
	Itrs           Interactors
	SessionManager s.MySessionManager
	i.OAuthHandler
	a.Annotator
}

func NewApp(itrs Interactors, sm s.MySessionManager, ip i.OAuthHandler, an a.Annotator) App {
	return App{itrs, sm, ip, an}
}

type Presenter struct {
}

func (p Presenter) Success(r createusr.Response) {
	fmt.Println("created initial admin user with id:", r.Id)
}
func (p Presenter) Error(err error) {
	if errors.Is(err, e.ErrDuplicate) {
		return
	}
	fmt.Println(fmt.Errorf("creating initial admin user: %w", err))
}

func MaybeCreateInitialAdmin(itr createusr.Interactor, email, password string) {
	initAdminUser := u.NewUser("admin", u.WithAdmin(true))
	itr.Execute(
		u.AppendUserToContext(context.Background(), initAdminUser),
		createusr.Request{Id: email, Password: &password, IsAdmin: true},
		Presenter{})

}
