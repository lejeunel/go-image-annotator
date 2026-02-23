package collections

import (
	g "datahub/generic"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
)

func RegisterCollectionsRoutes(api *huma.API, s *Service, URLBuilder g.APIURLBuilder) {
	collectionsController := &CollectionHTTPController{CollectionService: s}
	huma.Register(*api, huma.Operation{
		OperationID: "get-collection",
		Method:      http.MethodGet,
		Path:        URLBuilder.Build("collections/{name}"),
		Tags:        []string{"Collections"},
		Summary:     "Get collection",
		Description: "Get collection",
	}, collectionsController.Get)
	huma.Register(*api, huma.Operation{
		OperationID: "create-collection",
		Method:      http.MethodPost,
		Path:        URLBuilder.Build("collections"),
		Tags:        []string{"Collections"},
		Summary:     "Create a collection",
	}, collectionsController.Create)
	huma.Register(*api, huma.Operation{
		OperationID: "put-collection",
		Method:      http.MethodPut,
		Path:        URLBuilder.Build("collections/{name}"),
		Tags:        []string{"Collections"},
		Summary:     "Update a collection",
	}, collectionsController.Put)
	huma.Register(*api, huma.Operation{
		OperationID: "patch-collection",
		Method:      http.MethodPatch,
		Path:        URLBuilder.Build("collections/{name}"),
		Tags:        []string{"Collections"},
		Summary:     "Patch a collection",
	}, collectionsController.Patch)
	huma.Register(*api, huma.Operation{
		OperationID: "assign-profile",
		Method:      http.MethodPost,
		Path:        URLBuilder.Build("collections/{collectionname}/profile/{profilename}"),
		Tags:        []string{"Collections"},
		Summary:     "Assign profile",
	}, collectionsController.AssignProfile)
	huma.Register(*api, huma.Operation{
		OperationID: "unassign-profile",
		Method:      http.MethodDelete,
		Path:        URLBuilder.Build("collections/{collectionname}/profile"),
		Tags:        []string{"Collections"},
		Summary:     "Unassign profile",
	}, collectionsController.UnassignProfile)

}
