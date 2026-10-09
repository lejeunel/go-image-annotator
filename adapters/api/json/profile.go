package json

import (
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	s "github.com/lejeunel/go-image-annotator/shared"
)

type ProfilePresenter struct {
	Writer http.ResponseWriter
	ErrorPresenter
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

func (p ProfilePresenter) SuccessFindProfile(r pr.Profile) {
	s.WriteJSON(p.Writer, 200, MakeProfileResponse(r))
}

func (p ProfilePresenter) SuccessCreateProfile(r pr.Profile) {
	s.WriteJSON(p.Writer, 200, MakeProfileResponse(r))
}

func (p ProfilePresenter) SuccessDeleteProfile(string) {
	p.Writer.WriteHeader(http.StatusNoContent)
}

func (p ProfilePresenter) SuccessUpdateProfile(r pr.Profile) {
	s.WriteJSON(p.Writer, 200, MakeProfileResponse(r))
}

func NewProfilePresenter(w http.ResponseWriter, l slog.Logger) ProfilePresenter {
	return ProfilePresenter{Writer: w, ErrorPresenter: NewErrPresenter(w, l)}
}
