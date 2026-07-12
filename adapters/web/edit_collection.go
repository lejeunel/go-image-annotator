package web

import (
	"net/http"

	"github.com/lejeunel/go-image-annotator/use-cases/collection/update"
)

type EditCollectionPresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(update.Response) string
	HTMXErrorPresenter
}

func NewEditCollectionPresenter(w http.ResponseWriter) EditCollectionPresenter {
	task := "Updating collection"
	okMessageFunc := func(r update.Response) string {
		return "Successfully updated collection"
	}
	return EditCollectionPresenter{w, task, okMessageFunc, NewHTMXErrorPresenter(task, w)}
}

func (p EditCollectionPresenter) SuccessUpdateCollection(r update.Response) {
	payload, _ := NotifySuccessPayload(p.task, p.okMessageFunc(r))
	p.writer.Header().Set("HX-Trigger", string(payload))
	p.writer.WriteHeader(http.StatusOK)
}
func (s *Server) EditCollection(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	s.Collection.Update.Execute(r.Context(),
		update.Request{
			Name:           r.URL.Query().Get("name"),
			NewName:        r.FormValue("name"),
			NewDescription: r.FormValue("description"),
		},
		NewEditCollectionPresenter(w))
}
