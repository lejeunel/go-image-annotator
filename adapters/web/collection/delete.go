package collection

import (
	"fmt"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/delete"
)

type DeleteCollectionPresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(delete.Response) string
	htmx.ErrorPresenter
}

func NewDeleteCollectionPresenter(w http.ResponseWriter) DeleteCollectionPresenter {
	task := "Deleting collection"
	okMessageFunc := func(r delete.Response) string {
		return fmt.Sprintf("Successfully deleted collection %v", r.Name)
	}
	return DeleteCollectionPresenter{w, task, okMessageFunc, htmx.NewErrorPresenter(task, w)}
}

func (p DeleteCollectionPresenter) Success(r delete.Response) {
	payload, _ := htmx.NotifySuccessPayloadAndReload(p.task, p.okMessageFunc(r))
	p.writer.Header().Set("HX-Trigger", string(payload))
	p.writer.WriteHeader(http.StatusOK)
}
func (s *Server) Delete(w http.ResponseWriter, r *http.Request) {
	s.DeleteItr.Execute(r.Context(),
		r.URL.Query().Get("name"),
		NewDeleteCollectionPresenter(w))
}
