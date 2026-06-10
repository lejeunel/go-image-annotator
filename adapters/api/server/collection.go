package server

import (
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	presenter "github.com/lejeunel/go-image-annotator/adapters/api/json/collection"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/create"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/list"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/read"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/update"
)

func (s *Server) FindCollectionByName(w http.ResponseWriter, r *http.Request, name string) {
	s.Collection.Find.Execute(r.Context(), read.Request{Name: name},
		presenter.NewFindPresenter(w, s.Logger))
}
func (s *Server) CreateCollection(w http.ResponseWriter, r *http.Request) {
	body, ok := json.MustDecodeJSON[models.NewCollection](w, r)
	if !ok {
		return
	}

	s.Collection.Create.Execute(
		r.Context(),
		create.Request{Name: body.Name, Description: *body.Description},
		presenter.NewCreatePresenter(w, s.Logger))
}
func (s *Server) DeleteCollectionByName(w http.ResponseWriter, r *http.Request, name string) {
	s.Collection.Delete.Execute(
		r.Context(),
		delete.Request{Name: name},
		presenter.NewDeletePresenter(w, s.Logger))

}
func (s *Server) ListCollections(w http.ResponseWriter, r *http.Request, params ListCollectionsParams) {
	req := list.Request{Page: 1, PageSize: s.Collection.DefaultPageSize}
	if p := params.Page; p != nil {
		req.Page = *p
	}
	if p := params.PageSize; p != nil {
		req.PageSize = *p
	}
	s.Collection.List.Execute(r.Context(), req,
		presenter.NewListPresenter(w, s.Logger))

}

func (s *Server) UpdateCollectionByName(w http.ResponseWriter, r *http.Request, name string) {
	body, ok := json.MustDecodeJSON[models.UpdateCollection](w, r)
	if !ok {
		return
	}

	s.Collection.Update.Execute(r.Context(),
		update.Request{Name: name, NewName: body.Name, NewDescription: body.Description},
		presenter.NewUpdatePresenter(w, s.Logger))

}
