package profile

import (
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	s "github.com/lejeunel/go-image-annotator/shared"
)

type Find struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func MakeProfileResponse(p pr.Profile) models.Profile {
	labels := p.Labels
	if labels == nil {
		labels = []string{}
	}

	return models.Profile{
		Name:        p.Name,
		Description: p.Description,
		Group:       p.Group,
		Labels:      labels,
	}
}

func (p Find) SuccessFindProfile(r pr.Profile) {
	s.WriteJSON(p.Writer, 200, MakeProfileResponse(r))
}

func NewFindPresenter(w http.ResponseWriter, l slog.Logger) Find {
	return Find{Writer: w, ErrorPresenter: json.NewErrPresenter(w, l)}
}
