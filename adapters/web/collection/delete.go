package collection

import (
	"fmt"
	"net/http"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/delete"
)

type DeletePresenter struct {
	writer http.ResponseWriter
	b.RowURL
	task          string
	okMessageFunc func(delete.Response) string
	htmx.ErrorPresenter
}

func NewDeletePresenter(w http.ResponseWriter, u b.RowURL) DeletePresenter {
	task := "deleting collection"
	okMessageFunc := func(r delete.Response) string {
		return fmt.Sprintf("Successfully deleted collection %v", r.Name)
	}
	return DeletePresenter{w, u, task, okMessageFunc, htmx.NewErrorPresenter(task, w)}
}

func (p DeletePresenter) Success(r delete.Response) {
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, p.okMessageFunc(r))
}

func (p DeletePresenter) SuccessFindCollection(c clc.Collection) {
	b.RenderConfirmDeleteRow(len(listCollectionsFields),
		c.Name, "collection", p.Url, p.writer)
}
func (s *Server) Delete(w http.ResponseWriter, r *http.Request) {
	s.DeleteItr.Execute(r.Context(),
		r.URL.Query().Get(resourceUrlFieldName),
		NewDeletePresenter(w, s.RowURL))
}
