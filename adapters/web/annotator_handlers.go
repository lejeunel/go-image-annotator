package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	aw "github.com/lejeunel/go-image-annotator/adapters/web/annotator"
	"github.com/lejeunel/go-image-annotator/shared/html"
	"github.com/lejeunel/go-image-annotator/use-cases/annotate/remove"
	updlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/update-label"
)

func (s *Server) ViewImage(w http.ResponseWriter, r *http.Request) {
	view := aw.NewAnnotationView()
	s.annotator.Init(r.URL.Query().Get("id"), r.URL.Query().Get("collection"), view)
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
	view := aw.NewAnnotationView()
	s.annotator.Init(r.URL.Query().Get("id"), r.URL.Query().Get("collection"), view)
	view.RenderAnnotationList(w)
}
func (s *Server) GetAnnotationsAsJSON(w http.ResponseWriter, r *http.Request) {
	view := aw.NewAnnotationView()
	s.annotator.Init(r.URL.Query().Get("id"), r.URL.Query().Get("collection"), view)
	view.RenderAnnotations(w)
}
func (s *Server) DeleteAnnotation(w http.ResponseWriter, r *http.Request) {
	s.annotator.DeleteAnnotation(remove.Request{Id: r.URL.Query().Get("id")}, aw.NewAnnotationView())
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
	baseErr := fmt.Errorf("setting label")
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
