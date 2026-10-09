package json

import (
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	s "github.com/lejeunel/go-image-annotator/shared"
	"github.com/lejeunel/go-image-annotator/use-cases/group/create"
	"github.com/lejeunel/go-image-annotator/use-cases/group/update"
)

func MakeGroupResponse(group g.Group) models.Group {
	return models.Group{
		Name:        group.Name,
		Description: group.Description,
	}
}

type GroupPresenter struct {
	Writer http.ResponseWriter
	ErrorPresenter
}

func (p GroupPresenter) SuccessFindGroup(r g.Group) {
	s.WriteJSON(p.Writer, 200, MakeGroupResponse(r))
}

func (p GroupPresenter) SuccessCreateGroup(r create.Response) {
	s.WriteJSON(p.Writer, 200, models.Group{
		Name:        r.Name,
		Description: r.Description,
	})
}

func (p GroupPresenter) SuccessUpdateGroup(r update.Response) {
	s.WriteJSON(p.Writer, 200, models.Group{
		Name:        r.Name,
		Description: r.Description,
	})
}

func (p GroupPresenter) SuccessDeleteGroup(string) {
	p.Writer.WriteHeader(http.StatusNoContent)
}

func (p GroupPresenter) SuccessListGroups(r []g.Group) {
	groups := []models.Group{}
	for _, group := range r {
		groups = append(groups, MakeGroupResponse(group))
	}

	s.WriteJSON(p.Writer, 200, models.ListGroups{Groups: groups})
}

func NewGroupPresenter(w http.ResponseWriter, l slog.Logger) GroupPresenter {
	return GroupPresenter{Writer: w, ErrorPresenter: NewErrPresenter(w, l)}
}
