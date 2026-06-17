package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	aw "github.com/lejeunel/go-image-annotator/adapters/web/annotator"
	an "github.com/lejeunel/go-image-annotator/adapters/web/annotator/annotorious"
	ap "github.com/lejeunel/go-image-annotator/modules/annotator/presenters"
	assign_label "github.com/lejeunel/go-image-annotator/use-cases/annotate/assign-label"
	"github.com/lejeunel/go-image-annotator/use-cases/annotate/remove"
	updlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/update-label"
)

func (s *Server) ViewImage(w http.ResponseWriter, r *http.Request) {
	s.PageBuilder.SetUserIdentityFromContext(r.Context())
	view := aw.NewAnnotationView(s.PageBuilder)
	p := ap.NewAnnotationPagePresenter(view)
	s.Annotator.Init(r.Context(), r.URL.Query().Get("id"), r.URL.Query().Get("collection"), p, p, p)
	view.Render(w)
}
func (s *Server) MakeHTMLAnnotationPanel(w http.ResponseWriter, r *http.Request) {
	view := aw.NewAnnotationView(s.PageBuilder)
	p := ap.NewAnnotationPagePresenter(view)
	s.Annotator.Init(r.Context(), r.URL.Query().Get("id"), r.URL.Query().Get("collection"), p, p, p)
	view.RenderAnnotationList(w)
}
func (s *Server) SubmitLabel(w http.ResponseWriter, r *http.Request) {
	req := assign_label.Request{ImageId: r.URL.Query().Get("image_id"),
		Collection: r.URL.Query().Get("collection"), Label: r.URL.Query().Get("label")}
	p := ap.NewJSONPresenter()
	s.Annotator.AddLabel(r.Context(), req, &p)
	p.Write(w)
}
func (s *Server) SubmitPolygon(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := io.ReadAll(r.Body)
	var polyreq an.AnnotoriousPolygonRequest
	err := json.Unmarshal(bodyBytes, &polyreq)
	if err != nil {
		http.Error(w, fmt.Errorf("submit polygon: unmarshalling body: %w", err).Error(), http.StatusBadRequest)
		return
	}
	p := ap.NewJSONPresenter()
	s.Annotator.AddPolygon(r.Context(), an.ToAddPolygonRequest(polyreq), &p)
	p.Write(w)
}
func (s *Server) UpdatePolygon(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := io.ReadAll(r.Body)

	var polyreq an.AnnotoriousPolygonModel
	err := json.Unmarshal(bodyBytes, &polyreq)
	if err != nil {
		http.Error(w, fmt.Errorf("updating polygon: unmarshalling body: %w", err).Error(), http.StatusBadRequest)
		return
	}
	p := ap.NewJSONPresenter()
	s.Annotator.UpdatePolygon(r.Context(), an.ToUpdatePolygonRequest(polyreq), &p)
	p.Write(w)

}
func (s *Server) SubmitBox(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := io.ReadAll(r.Body)

	var boxreq an.AnnotoriousBoxRequest
	err := json.Unmarshal(bodyBytes, &boxreq)
	if err != nil {
		http.Error(w, fmt.Errorf("submit box: unmarshalling body: %w", err).Error(), http.StatusBadRequest)
		return
	}

	p := ap.NewJSONPresenter()
	s.Annotator.AddBox(r.Context(), an.ToAddBoxRequest(boxreq), &p)
	p.Write(w)
}
func (s *Server) UpdateBox(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := io.ReadAll(r.Body)

	var boxreq an.AnnotoriousBoxModel
	err := json.Unmarshal(bodyBytes, &boxreq)
	if err != nil {
		http.Error(w, fmt.Errorf("updating box: unmarshalling body: %w", err).Error(), http.StatusBadRequest)
		return
	}
	p := ap.NewJSONPresenter()
	s.Annotator.UpdateBox(r.Context(), an.ToUpdateBoxRequest(boxreq), &p)
	p.Write(w)
}
func (s *Server) DeleteAnnotation(w http.ResponseWriter, r *http.Request) {
	p := ap.NewJSONPresenter()
	s.Annotator.DeleteAnnotation(r.Context(), remove.Request{Id: r.URL.Query().Get("id")}, &p)
	p.Write(w)
}
func (s *Server) SetLabel(w http.ResponseWriter, r *http.Request) {
	errCtx := fmt.Errorf("setting label")
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, fmt.Errorf("%w: failed parsing url to get annotation id", errCtx).Error(), http.StatusBadRequest)
		return
	}
	label := r.URL.Query().Get("label")
	if label == "" {
		http.Error(w, fmt.Errorf("%w: failed parsing url to get label field", errCtx).Error(), http.StatusBadRequest)
		return
	}

	p := ap.NewJSONPresenter()
	s.Annotator.UpdateLabel(r.Context(), updlbl.Request{AnnotationId: id, Label: label}, &p)
	p.Write(w)
}
func (s *Server) GetRegionsAsJSON(w http.ResponseWriter, r *http.Request) {
	p := ap.NewJSONPresenter()
	s.Annotator.ReadImage(r.Context(), r.URL.Query().Get("id"), r.URL.Query().Get("collection"), &p)
	p.RenderRegionAnnotationsAsJSON(w)
}
