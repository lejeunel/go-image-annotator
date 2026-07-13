package web

import (
	"fmt"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/create"
	"net/http"
)

var createCollectionTargetDiv = "create-collection"

type CreateCollectionPresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(create.Response) string
	HTMXErrorPresenter
}

func NewCreateCollectionPresenter(w http.ResponseWriter) CreateCollectionPresenter {
	task := "Creating collection"
	okMessageFunc := func(r create.Response) string {
		return fmt.Sprintf("Successfully created collection %v", r.Name)
	}
	return CreateCollectionPresenter{w, task, okMessageFunc, NewHTMXErrorPresenter(task, w)}
}
func (p CreateCollectionPresenter) Success(r create.Response) {
	payload, _ := NotifySuccessPayloadAndReload(p.task, p.okMessageFunc(r))
	p.writer.Header().Set("HX-Trigger", string(payload))
	p.writer.WriteHeader(http.StatusOK)
}
func (s *Server) CreateCollection(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}
	s.Collection.Create.Execute(r.Context(), create.Request{Name: r.FormValue("name"),
		Description: r.FormValue("description")}, NewCreateCollectionPresenter(w))
}
