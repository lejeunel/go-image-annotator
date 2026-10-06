package profile

import (
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	s "github.com/lejeunel/go-image-annotator/shared"
)

type Update struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p Update) SuccessUpdateProfile(r pr.Profile) {
	s.WriteJSON(p.Writer, 200, MakeProfileResponse(r))
}

func NewUpdatePresenter(w http.ResponseWriter, l slog.Logger) Update {
	return Update{Writer: w, ErrorPresenter: json.NewErrPresenter(w, l)}
}
