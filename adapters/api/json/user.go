package json

import (
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	s "github.com/lejeunel/go-image-annotator/shared"
	"github.com/lejeunel/go-image-annotator/use-cases/user/create"
	"github.com/lejeunel/go-image-annotator/use-cases/user/list"
	upd "github.com/lejeunel/go-image-annotator/use-cases/user/update-privileges"
)

func MakeUserResponse(user u.User) models.User {
	return models.User{
		Id:     user.Id,
		Groups: user.GroupNames(),
		Roles:  user.RoleNames(),
	}
}

type UserPresenter struct {
	Writer http.ResponseWriter
	ErrorPresenter
}

func (p UserPresenter) SuccessCreateUser(r create.Response) {
	response := models.User{
		Id:     r.Id,
		Roles:  r.Roles,
		Groups: r.Groups,
	}

	s.WriteJSON(p.Writer, 200, response)
}

func (p UserPresenter) SuccessDeleteUser(id u.UserId) {
	p.Writer.WriteHeader(http.StatusNoContent)
}

func (p UserPresenter) SuccessUpdate(r upd.Response) {
	p.Writer.WriteHeader(http.StatusNoContent)
}

func (p UserPresenter) SuccessListUsers(r list.Response) {
	users := []models.User{}
	for _, user := range r.Users {
		users = append(users, MakeUserResponse(user))
	}

	s.WriteJSON(p.Writer, 200, models.ListUsers{
		Users:      users,
		Pagination: BuildPaginationResponse(r.Pagination),
	})
}

func NewUserPresenter(w http.ResponseWriter, l slog.Logger) UserPresenter {
	return UserPresenter{Writer: w, ErrorPresenter: NewErrPresenter(w, l)}
}
