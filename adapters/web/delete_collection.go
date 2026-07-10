package web

import (
	"fmt"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/delete"
	"net/http"
)

type DeleteCollectionPresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(delete.Response) string
	WebErrorPresenter
}

func NewDeleteCollectionPresenter(w http.ResponseWriter) DeleteCollectionPresenter {
	task := "Deleting collection"
	okMessageFunc := func(r delete.Response) string {
		return fmt.Sprintf("Successfully deleted collection with name %v", r.Name)
	}
	return DeleteCollectionPresenter{w, task, okMessageFunc, NewWebErrorPresenter(task, w)}
}

func (p DeleteCollectionPresenter) Success(r delete.Response) {
	payload, _ := NotifySuccessPayload(p.task, p.okMessageFunc(r))
	p.writer.Header().Set("HX-Trigger", string(payload))
	p.writer.WriteHeader(http.StatusOK)
}

func (s *Server) DeleteCollection(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	if name == "" {
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}
	s.Collection.Delete.Execute(r.Context(), name, NewDeleteCollectionPresenter(w))
}
