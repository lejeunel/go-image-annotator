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

	var description *string
	descriptionValue := r.FormValue(DescriptionFieldName)
	if descriptionValue != "" {
		description = &descriptionValue
	}
	s.Roles.Update.Execute(r.Context(),
		update.Request{
			Name:           name,
			NewName:        name,
			NewDescription: description,
			NewMethods:     methods,
		},
		NewEditPresenter(w, s.RowUrl))
}
