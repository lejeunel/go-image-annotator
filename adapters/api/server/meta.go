package server

import (
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	p "github.com/lejeunel/go-image-annotator/adapters/api/json/meta"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	"github.com/lejeunel/go-image-annotator/use-cases/metadata/add"
	"github.com/lejeunel/go-image-annotator/use-cases/metadata/delete"
)

func (s *Server) UpsertMetadata(
	w http.ResponseWriter,
	r *http.Request,
	collectionName string,
	imageId string,
) {
	body, ok := json.MustDecodeJSON[models.Meta](w, r)
	if !ok {
		return
	}

	req := add.Request{
		ImageId: imageId, Collection: collectionName,
		Key: body.Key, Value: body.Value,
	}

	s.Metadata.Add.Execute(r.Context(), req, p.NewPresenter(w, s.Logger))
}

func (s *Server) DeleteMetadata(
	w http.ResponseWriter,
	r *http.Request,
	collectionName string,
	imageId string,
	key string,
) {
	req := delete.Request{
		ImageId: imageId, Collection: collectionName,
		Key: key,
	}

	s.Metadata.Delete.Execute(r.Context(), req, p.NewPresenter(w, s.Logger))
}
