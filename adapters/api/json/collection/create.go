package collection

import (
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	c "github.com/lejeunel/go-image-annotator/entities/collection"
	s "github.com/lejeunel/go-image-annotator/shared"
)

type Create struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p Create) SuccessCreateCollection(r c.Collection) {
	s.WriteJSON(p.Writer, 200, MakeCollectionResponse(r))
}

func NewCreatePresenter(w http.ResponseWriter, l slog.Logger) Create {
	return Create{Writer: w, ErrorPresenter: json.NewErrPresenter(w, l)}
}
