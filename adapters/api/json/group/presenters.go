package group

import (
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
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

type Presenter struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p Presenter) SuccessFindGroup(r g.Group) {
	s.WriteJSON(p.Writer, 200, MakeGroupResponse(r))
}

func (p Presenter) Success(r create.Response) {
	s.WriteJSON(p.Writer, 200, models.Group{
		Name:        r.Name,
		Description: r.Description,
	})
}

func (p Presenter) SuccessUpdateGroup(r update.Response) {
	s.WriteJSON(p.Writer, 200, models.Group{
		Name:        r.Name,
		Description: r.Description,
	})
}

func (p Presenter) SuccessDeleteGroup(string) {
	p.Writer.WriteHeader(http.StatusNoContent)
}

func (p Presenter) SuccessListGroups(r []g.Group) {
	groups := []models.Group{}
	for _, group := range r {
		groups = append(groups, MakeGroupResponse(group))
	}

	s.WriteJSON(p.Writer, 200, models.ListGroups{Groups: groups})
}

func NewPresenter(w http.ResponseWriter, l slog.Logger) Presenter {
	return Presenter{Writer: w, ErrorPresenter: json.NewErrPresenter(w, l)}
}
