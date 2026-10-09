package json

import (
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	s "github.com/lejeunel/go-image-annotator/shared"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/clone"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/list"
)

type CollectionPresenter struct {
	Writer http.ResponseWriter
	ErrorPresenter
}

func MakeCollectionResponse(c clc.Collection) models.Collection {
	return models.Collection{
		Name:        c.Name,
		Description: c.Description,
		Group:       c.Group,
		Profile:     c.Profile,
	}
}

func (p CollectionPresenter) SuccessFindCollection(r clc.Collection) {
	s.WriteJSON(p.Writer, 200, MakeCollectionResponse(r))
}

func (p CollectionPresenter) SuccessSubmitCloneTask(r clone.Response) {
	response := models.TaskResponse{
		TaskId: r.Id, Issuer: r.Issuer, Type: r.Type,
	}

	s.WriteJSON(p.Writer, 200, response)
}

func (p CollectionPresenter) SuccessCreateCollection(r clc.Collection) {
	s.WriteJSON(p.Writer, 200, MakeCollectionResponse(r))
}

func (p CollectionPresenter) SuccessDeleteCollection(delete.Response) {
	p.Writer.WriteHeader(http.StatusNoContent)
}

func (p CollectionPresenter) SuccessUpdateCollection(r clc.Collection) {
	s.WriteJSON(p.Writer, 200, MakeCollectionResponse(r))
}

func (p CollectionPresenter) SuccessListCollections(r list.Response) {
	collections := []models.Collection{}
	for _, c := range r.Collections {
		collections = append(collections,
			MakeCollectionResponse(c))
	}

	response := models.ListCollections{
		Collections: collections,
		Pagination:  BuildPaginationResponse(r.Pagination),
	}

	s.WriteJSON(p.Writer, 200, response)
}

func NewCollectionPresenter(w http.ResponseWriter, l slog.Logger) CollectionPresenter {
	return CollectionPresenter{Writer: w, ErrorPresenter: NewErrPresenter(w, l)}
}
