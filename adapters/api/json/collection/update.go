package collection

import (
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	c "github.com/lejeunel/go-image-annotator/entities/collection"
	s "github.com/lejeunel/go-image-annotator/shared"
)

type Update struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p Update) SuccessUpdateCollection(r c.Collection) {
	s.WriteJSON(p.Writer, 200, MakeCollectionResponse(r))
}

func NewUpdatePresenter(w http.ResponseWriter, l slog.Logger) Update {
	return Update{Writer: w, ErrorPresenter: json.NewErrPresenter(w, l)}
}
