package server

import (
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	"github.com/lejeunel/go-image-annotator/adapters/api/json/annotate"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	addbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-bbox"
	addply "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-polygon"
	updlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/update-label"
)

func (s *Server) AddBoundingBox(w http.ResponseWriter, r *http.Request) {
	body, ok := json.MustDecodeJSON[models.AddBoundingBox](w, r)
	if !ok {
		return
	}
	req := addbox.Request{
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

func (s *Server) AddPolygon(w http.ResponseWriter, r *http.Request) {
	body, ok := json.MustDecodeJSON[models.AddPolygon](w, r)
	if !ok {
		return
	}
	req := addply.Request{
		ImageId:    body.ImageId,
		Collection: body.Collection,
		Label:      body.Label,
	}

	for _, point := range body.Points {
		req.Points.Append(point[0], point[1])
	}
	s.Annotation.AddPolygon.Execute(r.Context(), req, annotate.NewAnnotationPresenter(w, s.Logger))
}

func (s *Server) DeleteAnnotationById(w http.ResponseWriter, r *http.Request, id string) {
	s.Annotation.Delete.Execute(r.Context(), id, annotate.NewAnnotationPresenter(w, s.Logger))
}

func (s *Server) UpdateAnnotationById(
	w http.ResponseWriter,
	r *http.Request,
	id string,
	label string,
) {
	s.Annotation.UpdateLabel.Execute(
		r.Context(),
		updlbl.Request{Id: id, Label: label},
		annotate.NewAnnotationPresenter(w, s.Logger),
	)
}
