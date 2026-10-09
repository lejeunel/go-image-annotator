package json

import (
	"log/slog"
	"net/http"

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

type RolePresenter struct {
	Writer http.ResponseWriter
	ErrorPresenter
}

func (p RolePresenter) SuccessFindRole(r rl.Role) {
	s.WriteJSON(p.Writer, 200, MakeRoleResponse(r))
}

func (p RolePresenter) SuccessCreateRole(r rl.Role) {
	s.WriteJSON(p.Writer, 200, MakeRoleResponse(r))
}

func (p RolePresenter) SuccessUpdateRole(r rl.Role) {
	s.WriteJSON(p.Writer, 200, MakeRoleResponse(r))
}

func (p RolePresenter) SuccessDeleteRole(string) {
	p.Writer.WriteHeader(http.StatusNoContent)
}

func (p RolePresenter) SuccessListRoles(r []rl.Role) {
	roles := []models.Role{}
	for _, role := range r {
		roles = append(roles, MakeRoleResponse(role))
	}

	s.WriteJSON(p.Writer, 200, models.ListRoles{Roles: roles})
}

func NewRolePresenter(w http.ResponseWriter, l slog.Logger) RolePresenter {
	return RolePresenter{Writer: w, ErrorPresenter: NewErrPresenter(w, l)}
}
