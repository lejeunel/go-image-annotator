package web

import (
	"fmt"
	"net/http"

	rt "github.com/lejeunel/go-image-annotator/routes"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/delete"
)

type DeleteCollectionPresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(delete.Response) string
	HTMXErrorPresenter
}

func NewDeleteCollectionPresenter(w http.ResponseWriter) DeleteCollectionPresenter {
	task := "Deleting collection"
	okMessageFunc := func(r delete.Response) string {
		return fmt.Sprintf("Successfully deleted collection %v", r.Name)
	}
	return DeleteCollectionPresenter{w, task, okMessageFunc, NewHTMXErrorPresenter(task, w)}
}

func (p DeleteCollectionPresenter) Success(r delete.Response) {
	payload, _ := NotifySuccessPayload(p.task, p.okMessageFunc(r))
	p.writer.Header().Set("HX-Trigger", string(payload))
	p.writer.WriteHeader(http.StatusOK)
}
func (s *Server) ConfirmDeleteCollection(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	RenderConfirmDeleteRow(len(listCollectionsFields),
		name,
		"collection",
		rt.AppendValueToQueryArgs(rt.Collection, "name", name),
		w)
}
func (s *Server) DeleteCollection(w http.ResponseWriter, r *http.Request) {
	s.Collection.Delete.Execute(r.Context(),
		r.URL.Query().Get("name"),
		NewDeleteCollectionPresenter(w))
}
