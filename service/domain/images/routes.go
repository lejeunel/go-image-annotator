package images

import (
	g "datahub/generic"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterImagesRoutes(api *huma.API, imagesController *ImagesHTTPController,
	URLBuilder g.APIURLBuilder) {

	huma.Register(*api, huma.Operation{
		OperationID: "get-images",
		Method:      http.MethodGet,
		Path:        URLBuilder.Build("images"),
		Tags:        []string{"Images"},
		Summary:     "Get images",
	}, imagesController.GetImages)
	huma.Register(*api, huma.Operation{
		OperationID: "get-image-in-collection",
		Method:      http.MethodGet,
		Path:        URLBuilder.Build("collections/{collection_name}/images/{id}"),
		Tags:        []string{"Images"},
		Summary:     "Get one image in collection",
	}, imagesController.GetOne)
	huma.Register(*api, huma.Operation{
		OperationID: "delete-image",
		Method:      http.MethodDelete,
		Path:        URLBuilder.Build("collections/{collection_name}/images/{id}"),
		Tags:        []string{"Images"},
		Summary:     "Delete image from collection",
	}, imagesController.Delete)
	huma.Register(*api, huma.Operation{
		OperationID: "delete-all-images-in-collection",
		Method:      http.MethodDelete,
		Path:        URLBuilder.Build("collections/{id}/images"),
		Tags:        []string{"Images"},
		Summary:     "Delete a collection and all images it contains",
	}, imagesController.DeleteAll)
	huma.Register(*api, huma.Operation{
		OperationID: "update-annotation",
		Method:      http.MethodPut,
		Path:        URLBuilder.Build("annotations/{id}"),
		Tags:        []string{"Annotations"},
		Summary:     "Update annotation",
	}, imagesController.UpdateAnnotation)
	huma.Register(*api, huma.Operation{
		OperationID: "import-to-collection",
		Method:      http.MethodPost,
		Path:        URLBuilder.Build("collections/{name}/import"),
		Tags:        []string{"Collections"},
		Summary:     "Import images into collection",
	}, imagesController.Import)

}

func RegisterBaseImageRoutes(api *huma.API, imagesController *ImagesHTTPController, URLBuilder *g.APIURLBuilder) {

	huma.Register(*api, huma.Operation{
		OperationID: "get-image",
		Method:      http.MethodGet,
		Path:        URLBuilder.Build("base-image/{id}"),
		Tags:        []string{"Base-Image"},
		Summary:     "Get base-image",
	}, imagesController.GetBase)
	huma.Register(*api, huma.Operation{
		OperationID: "update-image",
		Method:      http.MethodPut,
		Path:        URLBuilder.Build("base-image/{id}"),
		Tags:        []string{"Base-Image"},
		Summary:     "Update base-image",
	}, imagesController.Update)
	huma.Register(*api, huma.Operation{
		OperationID: "patch-image",
		Method:      http.MethodPatch,
		Path:        URLBuilder.Build("base-image/{id}"),
		Tags:        []string{"Base-Image"},
		Summary:     "Patch base-image",
	}, imagesController.Patch)

}
