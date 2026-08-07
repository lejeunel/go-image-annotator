package policy

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	a "github.com/lejeunel/go-image-annotator/modules/authorizer"
	rt "github.com/lejeunel/go-image-annotator/routes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

func (s *Server) Route(r chi.Router, mws ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(mws...)
		r.Get(rt.AdminPoliciesUrl, s.Edit)
		r.Get(DefaultPolicyDownloadUrl, s.DownloadDefault)
		r.Post(SetPolicyFormUrl, s.Update)
	})
}

func (s *Server) Edit(w http.ResponseWriter, r *http.Request) {
	s.Page.SetUserIdentity(r.Context())
	s.Itrs.Read.Execute(r.Context(), NewViewPresenter(w, s.Page))
}

func (s *Server) Update(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	s.Itrs.Set.Execute(r.Context(), r.FormValue(PolicyFieldName),
		NewSetPresenter(w))
}

func (s *Server) DownloadDefault(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	if err := a.MarshalPolicies(a.DefaultPolicies, &buf); err != nil {
		http.Error(w,
			fmt.Errorf("generating default yaml policies: %v: %w", err, e.ErrInternal).Error(),
			http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/x-yaml")
	w.Header().Set("Content-Disposition", `attachment; filename="`+a.DefaultPolicyFileName+`"`)
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}
