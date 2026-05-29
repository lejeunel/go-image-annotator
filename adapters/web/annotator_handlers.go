package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	aw "github.com/lejeunel/go-image-annotator/adapters/web/annotator"
	an "github.com/lejeunel/go-image-annotator/entities/annotation"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/lejeunel/go-image-annotator/shared/html"
	"github.com/lejeunel/go-image-annotator/use-cases/annotate/remove"
	updlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/update-label"
)

func ParseAnnotationIdFromURL(u *url.URL) (*an.AnnotationId, error) {
	baseErr := "parsing url for annotation id"
	idStr := u.Query().Get("id")
	if idStr == "" {
		return nil, fmt.Errorf("%v: extracting id: %w", baseErr, e.ErrURLParsing)
	}
	id, err := an.NewAnnotationIdFromString(idStr)
	if err != nil {
		return nil, fmt.Errorf("%v: validating id (%v): %w", baseErr, idStr, e.ErrValidation)
	}
	return id, nil
}

func ParseImageIdAndCollectionFromURL(u *url.URL) (*aw.Request, error) {
	baseErr := "parsing url"
	req := aw.Request{}
	imageIdStr := u.Query().Get("id")
	if imageIdStr == "" {
		return nil, fmt.Errorf("%v: extracting id: %w", baseErr, e.ErrURLParsing)
	}
	req.ImageId = imageIdStr

	collection := u.Query().Get("collection")
	if collection == "" {
		return nil, fmt.Errorf("%v: collection (%v): %w", baseErr, collection, e.ErrURLParsing)
	}
	req.Collection = collection
	return &req, nil
}
func (s *Server) ViewImage(w http.ResponseWriter, r *http.Request) {

	req, err := ParseImageIdAndCollectionFromURL(r.URL)
	if err != nil {
		html.NewPageBuilder().SetError(err).Render(w)
		return
	}

	view := aw.NewAnnotationView()
	s.annotator.Init(req.ImageId, req.Collection, view)
	view.RenderAll(w)
}
func (s *Server) SubmitBox(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := io.ReadAll(r.Body)

	var boxreq aw.AnnotoriousBoxRequest
	err := json.Unmarshal(bodyBytes, &boxreq)
	if err != nil {
		http.Error(w, fmt.Errorf("submit box: unmarshalling body: %w", err).Error(), http.StatusBadRequest)
		return
	}

	req, err := aw.ToAddBoxRequest(boxreq)
	if err != nil {
		http.Error(w, fmt.Errorf("submit box: converting box: %w", err).Error(), http.StatusBadRequest)
		return
	}
	s.annotator.AddBox(*req, aw.NewAnnotationView())
}
func (s *Server) MakeHTMLAnnotationPanel(w http.ResponseWriter, r *http.Request) {
	req, err := ParseImageIdAndCollectionFromURL(r.URL)
	if err != nil {
		html.NewPageBuilder().SetError(err).Render(w)
		return
	}
	view := aw.NewAnnotationView()
	s.annotator.Init(req.ImageId, req.Collection, view)
	view.RenderAnnotationList(w)
}
func (s *Server) GetAnnotationsAsJSON(w http.ResponseWriter, r *http.Request) {
	req, err := ParseImageIdAndCollectionFromURL(r.URL)
	if err != nil {
		html.NewPageBuilder().SetError(err).Render(w)
		return
	}
	view := aw.NewAnnotationView()
	s.annotator.Init(req.ImageId, req.Collection, view)
	view.RenderAnnotations(w)
}
func (s *Server) DeleteAnnotation(w http.ResponseWriter, r *http.Request) {
	id, err := ParseAnnotationIdFromURL(r.URL)
	if err != nil {
		html.NewPageBuilder().SetError(err).Render(w)
		return
	}
	s.annotator.DeleteAnnotation(remove.Request{Id: *id}, aw.NewAnnotationView())
}
func (s *Server) UpdateBox(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := io.ReadAll(r.Body)

	var boxreq aw.AnnotoriousBoxModel
	err := json.Unmarshal(bodyBytes, &boxreq)
	if err != nil {
		http.Error(w, fmt.Errorf("submit box: unmarshalling body: %w", err).Error(), http.StatusBadRequest)
		return
	}

	req, err := aw.ToUpdateBoxRequest(boxreq)
	if err != nil {
		http.Error(w, fmt.Errorf("submit box: converting box: %w", err).Error(), http.StatusBadRequest)
		return
	}
	s.annotator.UpdateBox(*req, aw.NewAnnotationView())
}

func (s *Server) SetLabel(w http.ResponseWriter, r *http.Request) {
	baseErr := "setting label"
	id := r.URL.Query().Get("id")
	if id == "" {
		html.NewPageBuilder().SetError(fmt.Errorf("%w: failed parsing url to get annotation id", baseErr)).Render(w)
		return
	}
	label := r.URL.Query().Get("label")
	if label == "" {
		html.NewPageBuilder().SetError(fmt.Errorf("%w: failed parsing url to get label field", baseErr)).Render(w)
		return
	}

	s.annotator.UpdateLabel(updlbl.Request{AnnotationId: id, Label: label},
		aw.NewAnnotationView())
}
