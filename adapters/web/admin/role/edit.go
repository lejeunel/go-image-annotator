package role

import (
	"net/http"

	"github.com/lejeunel/go-image-annotator/use-cases/role/update"
)

func (s *Server) Edit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	name := r.URL.Query().Get(resourceUrlFieldName)
	methods := r.Form[MethodsFieldName]

	req := update.Request{
		Name:       name,
		NewName:    name,
		NewMethods: methods,
	}
	d := r.FormValue(DescriptionFieldName)
	if d != "" {
		req.NewDescription = &d
	}
	s.Roles.Update.Execute(r.Context(),
		req,
		NewEditPresenter(w, s.RowUrl))
}
