package user

import (
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
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

type Presenter struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p Presenter) SuccessCreateUser(r create.Response) {
	response := models.User{
		Id:     r.Id,
		Roles:  r.Roles,
		Groups: r.Groups,
	}

	s.WriteJSON(p.Writer, 200, response)
}

func (p Presenter) SuccessDeleteUser(id u.UserId) {
	p.Writer.WriteHeader(http.StatusNoContent)
}

func (p Presenter) SuccessUpdate(r upd.Response) {
	p.Writer.WriteHeader(http.StatusNoContent)
}

func (p Presenter) SuccessListUsers(r list.Response) {
	users := []models.User{}
	for _, user := range r.Users {
		users = append(users, MakeUserResponse(user))
	}

	s.WriteJSON(p.Writer, 200, models.ListUsers{
		Users:      users,
		Pagination: json.BuildPaginationResponse(r.Pagination),
	})
}

func NewPresenter(w http.ResponseWriter, l slog.Logger) Presenter {
	return Presenter{Writer: w, ErrorPresenter: json.NewErrPresenter(w, l)}
}
