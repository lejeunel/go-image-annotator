package collection

import (
	"fmt"
	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/create"
	"net/http"
)

var createCollectionTargetDiv = "create-collection"

type CreateCollectionPresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(create.Response) string
	htmx.ErrorPresenter
}

func NewCreateCollectionPresenter(w http.ResponseWriter) CreateCollectionPresenter {
	task := "Creating collection"
	okMessageFunc := func(r create.Response) string {
		return fmt.Sprintf("Successfully created collection %v", r.Name)
	}
	return CreateCollectionPresenter{w, task, okMessageFunc, htmx.NewErrorPresenter(task, w)}
}
func (p CreateCollectionPresenter) Success(r create.Response) {
	payload, _ := htmx.NotifySuccessPayloadAndReload(p.task, p.okMessageFunc(r))
	p.writer.Header().Set("HX-Trigger", string(payload))
	p.writer.WriteHeader(http.StatusOK)
}
func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}
	s.CreateItr.Execute(r.Context(), create.Request{Name: r.FormValue("name"),
		Description: r.FormValue("description")}, NewCreateCollectionPresenter(w))
}
