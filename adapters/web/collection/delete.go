package collection

import (
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
		return MakeNewTaskMessage()
	}
	return DeletePresenter{w, u, task, okMessageFunc, htmx.NewErrorPresenter(task, w)}
}

func (p DeletePresenter) SuccessDeleteCollection(r delete.Response) {
	htmx.NotifySuccessPayload(p.writer, p.task, p.okMessageFunc(r))
}

func (p DeletePresenter) SuccessFindCollection(c clc.Collection) {
	b.RenderConfirmDeleteRow(len(listCollectionsFields),
		c.Name, "collection", p.Url, p.writer)
}

func (s *Server) Delete(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get(resourceUrlFieldName)
	s.DeleteItr.Execute(r.Context(), name,
		NewDeletePresenter(w, s.RowURL))
}
