package server

import (
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/clone"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/create"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/update"
)

func (s *Server) FindCollectionByName(w http.ResponseWriter, r *http.Request, name string) {
	s.Collection.Find.Execute(r.Context(), name,
		json.NewCollectionPresenter(w, s.Logger))
}

func (s *Server) CreateCollection(w http.ResponseWriter, r *http.Request) {
	body, ok := json.MustDecodeJSON[models.NewCollection](w, r)
	if !ok {
		return
	}

	s.Collection.Create.Execute(
		r.Context(),
		create.Request{Name: body.Name, Description: body.Description},
		json.NewCollectionPresenter(w, s.Logger))
}

func (s *Server) DeleteCollection(w http.ResponseWriter, r *http.Request, name string) {
	s.Collection.Delete.Execute(r.Context(), name, json.NewCollectionPresenter(w, s.Logger))
}

func (s *Server) ListCollections(
	w http.ResponseWriter,
	r *http.Request,
	params ListCollectionsParams,
) {
	s.Collection.List.Execute(
		r.Context(),
		pa.NewPaginationParamsFromOptional(
			params.PageSize,
			params.Page,
			s.Collection.DefaultPageSize,
		),
		json.NewCollectionPresenter(w, s.Logger),
	)
}

func (s *Server) UpdateCollection(w http.ResponseWriter, r *http.Request, name string) {
	body, ok := json.MustDecodeJSON[models.UpdateCollection](w, r)
	if !ok {
		return
	}

	s.Collection.Update.Execute(r.Context(),
		update.Request{
			Name: name, NewName: body.Name, NewDescription: body.Description,
			NewGroup: body.Group, NewProfile: body.Profile,
		},
		json.NewCollectionPresenter(w, s.Logger))
}

func (s *Server) CloneCollection(w http.ResponseWriter, r *http.Request) {
	body, ok := json.MustDecodeJSON[models.CloneCollection](w, r)
	if !ok {
		return
	}

	req := clone.Request{
		Source: body.Source, Destination: body.Destination,
		DestinationGroup: body.Group, Deep: body.Deep,
	}
	s.Collection.Clone.Execute(r.Context(), req,
		json.NewCollectionPresenter(w, s.Logger))
}
