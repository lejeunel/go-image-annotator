package image

import (
	"fmt"
	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	"github.com/lejeunel/go-image-annotator/use-cases/image/delete"
	"net/http"
)

type DeleteImagePresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(delete.Response) string
	htmx.ErrorPresenter
}

func NewDeleteImagePresenter(w http.ResponseWriter) DeleteImagePresenter {
	task := "Deleting image"
	okMessageFunc := func(r delete.Response) string {
		return fmt.Sprintf("Successfully deleted image %v from collection %v", r.ImageId, r.Collection)
	}
	return DeleteImagePresenter{w, task, okMessageFunc, htmx.NewErrorPresenter(task, w)}
}

func (p DeleteImagePresenter) SuccessDeleteImage(r delete.Response) {
	payload, _ := htmx.NotifySuccessPayloadAndReload(p.task, p.okMessageFunc(r))
	p.writer.Header().Set("HX-Trigger", string(payload))
	p.writer.WriteHeader(http.StatusOK)
}

func (s *Server) Delete(w http.ResponseWriter, r *http.Request) {
	s.DeleteItr.Execute(r.Context(),
		delete.Request{ImageId: r.URL.Query().Get("id"), Collection: r.URL.Query().Get("collection")},
		NewDeleteImagePresenter(w))

}
