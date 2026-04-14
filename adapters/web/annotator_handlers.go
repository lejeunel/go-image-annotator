package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	aw "github.com/lejeunel/go-image-annotator-v2/adapters/web/annotator"
	an "github.com/lejeunel/go-image-annotator-v2/application/annotator"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	e "github.com/lejeunel/go-image-annotator-v2/shared/errors"
	"github.com/lejeunel/go-image-annotator-v2/shared/html"
	addbox "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/add-bbox"
)

func ParseURL(u *url.URL) (*Request, error) {
	baseErr := "parsing url"
	req := Request{}
	imageIdStr := u.Query().Get("id")
	if imageIdStr == "" {
		return nil, fmt.Errorf("%v: extracting id: %w", baseErr, e.ErrURLParsing)
	}
	imageId, err := im.NewImageIdFromString(imageIdStr)
	if err != nil {
		return nil, fmt.Errorf("%v: validating id (%v): %w", baseErr, imageIdStr, e.ErrValidation)
	}
	req.ImageId = imageId

	collection := u.Query().Get("collection")
	if collection == "" {
		return nil, fmt.Errorf("%v: collection (%v): %w", baseErr, collection, e.ErrURLParsing)
	}
	req.Collection = collection
	return &req, nil
}

func (s *Server) ViewImage(w http.ResponseWriter, r *http.Request) {

	req, err := ParseURL(r.URL)
	if err != nil {
		html.NewPageBuilder().SetError(err).Render(w)
		return
	}
	view := aw.NewAnnotationView()
	presenter := an.NewAnnotatorPresenter(view)
	s.annotator.Start(req.ImageId, req.Collection, *presenter)
	view.Render(w)
}

func convertBoxRequest(r BoxRequest) (*addbox.Request, error) {
	xc := r.Annotation.Selector.Geometry.XTopLeft + r.Annotation.Selector.Geometry.W/2
	yc := r.Annotation.Selector.Geometry.YTopLeft + r.Annotation.Selector.Geometry.H/2
	width := r.Annotation.Selector.Geometry.W
	height := r.Annotation.Selector.Geometry.H

	imageId, err := im.NewImageIdFromString(r.ImageId)
	if err != nil {
		return nil, fmt.Errorf("submitting annotation: validating imageId: %w", err)
	}

	return &addbox.Request{ImageId: imageId, Collection: r.Collection,
		Label: r.Label, Xc: float32(xc), Yc: float32(yc), Width: float32(width), Height: float32(height)}, nil

}

func (s *Server) SubmitBox(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := io.ReadAll(r.Body)

	var boxreq BoxRequest
	err := json.Unmarshal(bodyBytes, &boxreq)
	if err != nil {
		html.NewPageBuilder().SetError(fmt.Errorf("submitting annotation: reading json payload: %w", err)).Render(w)
		return
	}
	view := aw.NewAnnotationView()
	presenter := an.NewAnnotatorPresenter(view)

	req, err := convertBoxRequest(boxreq)
	if err != nil {
		html.NewPageBuilder().SetError(fmt.Errorf("submitting annotation: %w", err)).Render(w)
		return
	}

	s.annotator.AddBox(*req, presenter)
	view.Render(w)
}
