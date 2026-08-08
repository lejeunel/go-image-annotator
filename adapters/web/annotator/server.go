package annotator

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	ap "github.com/lejeunel/go-image-annotator/adapters/web/annotator/presenters"
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	a "github.com/lejeunel/go-image-annotator/modules/annotator"
	s "github.com/lejeunel/go-image-annotator/shared/session"
	assign_label "github.com/lejeunel/go-image-annotator/use-cases/annotate/assign-label"
	"github.com/lejeunel/go-image-annotator/use-cases/annotate/remove"
	updlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/update-label"
)

type Server struct {
	b.PageBuilder
	a.Annotator
	s.SessionManager
}

func NewServer(
	annotator a.Annotator,
	pageBuilder b.PageBuilder,
	sessionManager s.SessionManager,
) *Server {
	return &Server{
		Annotator:      annotator,
		SessionManager: sessionManager,
		PageBuilder:    *pageBuilder.SetHTMLTitle("Annotate"),
	}
}

func (s *Server) AnnotateImage(w http.ResponseWriter, r *http.Request) {
	s.PageBuilder.SetUserIdentity(r.Context())
	view := NewAnnotationView(s.PageBuilder)
	p := ap.NewAnnotationPagePresenter(ap.NewCyclicColorizer(ap.Palette))
	p.SetView(view)
	s.Annotator.Init(r.Context(), r.URL.Query().Get("id"), r.URL.Query().Get("collection"), p, p, p)
	view.Render(w)
}

func (s *Server) MakeAnnotationPanel(w http.ResponseWriter, r *http.Request) {
	view := NewAnnotationView(s.PageBuilder)
	p := ap.NewAnnotationPagePresenter(ap.NewCyclicColorizer(ap.Palette))
	p.SetView(view)
	s.Annotator.Init(r.Context(), r.URL.Query().Get("id"), r.URL.Query().Get("collection"), p, p, p)
	view.RenderAnnotationList(w)
}

func (s *Server) SubmitLabel(w http.ResponseWriter, r *http.Request) {
	req := assign_label.Request{
		ImageId:    r.URL.Query().Get("image_id"),
		Collection: r.URL.Query().Get("collection"), Label: r.URL.Query().Get("label"),
	}

	p := ap.NewAnnotoriousPresenter(w)
	s.Annotator.AddLabel.Execute(r.Context(), req, &p)
}

func (s *Server) SubmitPolygon(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := io.ReadAll(r.Body)
	var polyreq ap.AnnotoriousPolygonRequest
	err := json.Unmarshal(bodyBytes, &polyreq)
	if err != nil {
		http.Error(
			w,
			fmt.Errorf("submit polygon: unmarshalling body: %w", err).Error(),
			http.StatusBadRequest,
		)
		return
	}
	p := ap.NewAnnotoriousPresenter(w)
	s.Annotator.AddPolygon.Execute(r.Context(), ap.ToAddPolygonRequest(polyreq), &p)
}

func (s *Server) UpdatePolygon(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := io.ReadAll(r.Body)

	var polyreq ap.AnnotoriousPolygonModel
	err := json.Unmarshal(bodyBytes, &polyreq)
	if err != nil {
		http.Error(
			w,
			fmt.Errorf("updating polygon: unmarshalling body: %w", err).Error(),
			http.StatusBadRequest,
		)
		return
	}
	p := ap.NewAnnotoriousPresenter(w)
	s.Annotator.UpdatePolygon.Execute(r.Context(), ap.ToUpdatePolygonRequest(polyreq), &p)
}

func (s *Server) SubmitBox(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := io.ReadAll(r.Body)

	var boxreq ap.AnnotoriousBoxRequest
	err := json.Unmarshal(bodyBytes, &boxreq)
	if err != nil {
		http.Error(
			w,
			fmt.Errorf("submit box: unmarshalling body: %w", err).Error(),
			http.StatusBadRequest,
		)
		return
	}

	p := ap.NewAnnotoriousPresenter(w)
	s.Annotator.AddBox.Execute(r.Context(), ap.ToAddBoxRequest(boxreq), &p)
}

func (s *Server) UpdateBox(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := io.ReadAll(r.Body)

	var boxreq ap.AnnotoriousBoxModel
	err := json.Unmarshal(bodyBytes, &boxreq)
	if err != nil {
		http.Error(
			w,
			fmt.Errorf("updating box: unmarshalling body: %w", err).Error(),
			http.StatusBadRequest,
		)
		return
	}
	p := ap.NewAnnotoriousPresenter(w)
	s.Annotator.UpdateBox.Execute(r.Context(), ap.ToUpdateBoxRequest(boxreq), &p)
}

func (s *Server) DeleteAnnotation(w http.ResponseWriter, r *http.Request) {
	p := ap.NewAnnotoriousPresenter(w)
	s.Annotator.DeleteAnnotation.Execute(
		r.Context(),
		remove.Request{Id: r.URL.Query().Get("id")},
		&p,
	)
}

func (s *Server) SetLabel(w http.ResponseWriter, r *http.Request) {
	errCtx := fmt.Errorf("setting label")
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(
			w,
			fmt.Errorf("%w: failed parsing url to get annotation id", errCtx).Error(),
			http.StatusBadRequest,
		)
		return
	}
	label := r.URL.Query().Get("label")
	if label == "" {
		http.Error(
			w,
			fmt.Errorf("%w: failed parsing url to get label field", errCtx).Error(),
			http.StatusBadRequest,
		)
		return
	}
	p := ap.NewAnnotoriousPresenter(w)
	s.Annotator.UpdateLabel.Execute(r.Context(), updlbl.Request{AnnotationId: id, Label: label}, &p)
}

func (s *Server) GetRegionsAsJSON(w http.ResponseWriter, r *http.Request) {
	p := ap.NewAnnotoriousPresenter(w)
	s.Annotator.ReadImage(r.URL.Query().Get("id"), r.URL.Query().Get("collection"), &p)
	p.RenderRegionAnnotationsAsJSON(w)
}
