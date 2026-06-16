package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	aw "github.com/lejeunel/go-image-annotator/adapters/web/annotator"
	an "github.com/lejeunel/go-image-annotator/adapters/web/annotator/annotorious"
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	assign_label "github.com/lejeunel/go-image-annotator/use-cases/annotate/assign-label"
	"github.com/lejeunel/go-image-annotator/use-cases/annotate/remove"
	updlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/update-label"
)

func (s *Server) ViewImage(w http.ResponseWriter, r *http.Request) {
	s.PageBuilder.SetUserIdentityFromContext(r.Context())
	view := aw.NewAnnotationView(s.PageBuilder)
	s.Annotator.Init(r.Context(), r.URL.Query().Get("id"), r.URL.Query().Get("collection"), view)
	view.RenderAll(w)
}
func (s *Server) SubmitLabel(w http.ResponseWriter, r *http.Request) {
	req := assign_label.Request{ImageId: r.URL.Query().Get("image_id"),
		Collection: r.URL.Query().Get("collection"), Label: r.URL.Query().Get("label")}
	s.Annotator.AddLabel(r.Context(), req, aw.NewAnnotationView(s.PageBuilder))
}
func (s *Server) SubmitPolygon(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := io.ReadAll(r.Body)
	var polyreq an.AnnotoriousPolygonRequest
	err := json.Unmarshal(bodyBytes, &polyreq)
	if err != nil {
		http.Error(w, fmt.Errorf("submit polygon: unmarshalling body: %w", err).Error(), http.StatusBadRequest)
		return
	}
}
func (s *Server) SubmitBox(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := io.ReadAll(r.Body)

	var boxreq an.AnnotoriousBoxRequest
	err := json.Unmarshal(bodyBytes, &boxreq)
	if err != nil {
		http.Error(w, fmt.Errorf("submit box: unmarshalling body: %w", err).Error(), http.StatusBadRequest)
		return
	}

	req, err := an.ToAddBoxRequest(boxreq)
	if err != nil {
		http.Error(w, fmt.Errorf("submit box: converting box: %w", err).Error(), http.StatusBadRequest)
		return
	}
	s.Annotator.AddBox(r.Context(), *req, aw.NewAnnotationView(s.PageBuilder))
}
func (s *Server) MakeHTMLAnnotationPanel(w http.ResponseWriter, r *http.Request) {
	view := aw.NewAnnotationView(s.PageBuilder)
	s.Annotator.Init(r.Context(), r.URL.Query().Get("id"), r.URL.Query().Get("collection"), view)
	view.RenderAnnotationList(w)
}
func (s *Server) GetAnnotationsAsJSON(w http.ResponseWriter, r *http.Request) {
	view := aw.NewAnnotationView(s.PageBuilder)
	s.Annotator.Init(r.Context(), r.URL.Query().Get("id"), r.URL.Query().Get("collection"), view)
	view.RenderAnnotations(w)
}
func (s *Server) DeleteAnnotation(w http.ResponseWriter, r *http.Request) {
	s.Annotator.DeleteAnnotation(r.Context(), remove.Request{Id: r.URL.Query().Get("id")},
		aw.NewAnnotationView(s.PageBuilder))
}
func (s *Server) UpdateBox(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := io.ReadAll(r.Body)

	var boxreq an.AnnotoriousBoxModel
	err := json.Unmarshal(bodyBytes, &boxreq)
	if err != nil {
		http.Error(w, fmt.Errorf("submit box: unmarshalling body: %w", err).Error(), http.StatusBadRequest)
		return
	}
	req, err := an.ToUpdateBoxRequest(boxreq)
	if err != nil {
		http.Error(w, fmt.Errorf("submit box: converting box: %w", err).Error(), http.StatusBadRequest)
		return
	}
	s.Annotator.UpdateBox(r.Context(), *req, aw.NewAnnotationView(s.PageBuilder))
}
func (s *Server) SetLabel(w http.ResponseWriter, r *http.Request) {
	errCtx := fmt.Errorf("setting label")
	id := r.URL.Query().Get("id")
	if id == "" {
		b.NewPageBuilder(s.APIPath).SetError(fmt.Errorf("%w: failed parsing url to get annotation id", errCtx)).Render(w)
		return
	}
	label := r.URL.Query().Get("label")
	if label == "" {
		b.NewPageBuilder(s.APIPath).SetError(fmt.Errorf("%w: failed parsing url to get label field", errCtx)).Render(w)
		return
	}

	s.Annotator.UpdateLabel(r.Context(), updlbl.Request{AnnotationId: id, Label: label},
		aw.NewAnnotationView(s.PageBuilder))
}
