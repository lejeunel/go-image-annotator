package web

import (
	"fmt"
	"net/http"

	bd "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	rt "github.com/lejeunel/go-image-annotator/routes"
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
		return fmt.Sprintf("Successfully updated collection %v", r.Name)
	}
	return EditCollectionPresenter{w, task, okMessageFunc, NewHTMXErrorPresenter(task, w)}
}

func (p EditCollectionPresenter) SuccessUpdateCollection(r update.Response) {
	payload, _ := NotifySuccessPayload(p.task, p.okMessageFunc(r))
	p.writer.Header().Set("HX-Trigger", string(payload))
	p.writer.WriteHeader(http.StatusOK)
}
func (s *Server) EditCollectionForm(w http.ResponseWriter, r *http.Request) {
	b := bd.NewHTMXInlineFormBuilder(len(listCollectionsFields), rt.Collection, bd.HTMXPutMethod)
	b.AddTextField("name", "Name", "name", bd.WithRequired(), bd.WithDefault(r.URL.Query().Get("name")))
	b.AddTextField("description", "Description", "description", bd.WithDefault(r.URL.Query().Get("description")))
	b.Render(w)
}
func (s *Server) EditCollection(w http.ResponseWriter, r *http.Request) {
	s.Collection.Update.Execute(r.Context(),
		update.Request{
			r.URL.Query().Get("name"),
			r.URL.Query().Get("name"),
			r.URL.Query().Get("description"),
		},
		NewEditCollectionPresenter(w))
}
