package profile

import (
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	s "github.com/lejeunel/go-image-annotator/shared"
)

type Create struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p Create) SuccessCreateProfile(r pr.Profile) {
	s.WriteJSON(p.Writer, 200, MakeProfileResponse(r))
}

func NewCreatePresenter(w http.ResponseWriter, l slog.Logger) Create {
	return Create{Writer: w, ErrorPresenter: json.NewErrPresenter(w, l)}
}
