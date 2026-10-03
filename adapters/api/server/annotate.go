package server

import (
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	"github.com/lejeunel/go-image-annotator/adapters/api/json/annotate"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	add "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-bbox"
)

func (s *Server) AddBoundingBox(w http.ResponseWriter, r *http.Request) {
	body, ok := json.MustDecodeJSON[models.AddBoundingBox](w, r)
	if !ok {
		return
	}
	req := add.Request{
		ImageId:    body.ImageId,
		Collection: body.Collection,
		Label:      body.Label,
		Xc:         body.Xc,
		Yc:         body.Yc,
		Width:      body.Width,
		Height:     body.Height,
		Angle:      body.Angle,
	}
	s.Annotation.AddBox.Execute(r.Context(), req, annotate.NewAnnotationPresenter(w, s.Logger))
}
