package collection

import (
	"net/http"

	"github.com/lejeunel/go-image-annotator/use-cases/collection/update"
)

func (s *Server) Edit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	s.UpdateItr.Execute(r.Context(),
		update.Request{
			Name:           r.URL.Query().Get("name"),
			NewName:        r.FormValue("name"),
			NewDescription: r.FormValue("description"),
		},
		NewEditPresenter(w, s.RowURL))
}
