package role

import (
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	rl "github.com/lejeunel/go-image-annotator/entities/role"
	s "github.com/lejeunel/go-image-annotator/shared"
)

func MakeRoleResponse(role rl.Role) models.Role {
	methods := role.Methods
	if methods == nil {
		methods = []string{}
	}

	return models.Role{
		Name:        role.Name,
		Description: role.Description,
		Methods:     methods,
	}
}

type Presenter struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p Presenter) SuccessFindRole(r rl.Role) {
	s.WriteJSON(p.Writer, 200, MakeRoleResponse(r))
}

func (p Presenter) SuccessCreateRole(r rl.Role) {
	s.WriteJSON(p.Writer, 200, MakeRoleResponse(r))
}

func (p Presenter) SuccessUpdateRole(r rl.Role) {
	s.WriteJSON(p.Writer, 200, MakeRoleResponse(r))
}

func (p Presenter) SuccessDeleteRole(string) {
	p.Writer.WriteHeader(http.StatusNoContent)
}

func (p Presenter) SuccessListRoles(r []rl.Role) {
	roles := []models.Role{}
	for _, role := range r {
		roles = append(roles, MakeRoleResponse(role))
	}

	s.WriteJSON(p.Writer, 200, models.ListRoles{Roles: roles})
}

func NewPresenter(w http.ResponseWriter, l slog.Logger) Presenter {
	return Presenter{Writer: w, ErrorPresenter: json.NewErrPresenter(w, l)}
}
