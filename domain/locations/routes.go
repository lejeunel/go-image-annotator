package locations

import (
	g "datahub/generic"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
)

func RegisterLocationsRoutes(api *huma.API, l *Service, URLBuilder g.APIURLBuilder) {
	locationsControllers := &LocationHTTPController{LocationService: l}
	huma.Register(*api, huma.Operation{
		OperationID: "get-site",
		Method:      http.MethodGet,
		Path:        URLBuilder.Build("sites/{name}"),
		Tags:        []string{"Locations"},
		Summary:     "Get site",
	}, locationsControllers.GetSite)
	huma.Register(*api, huma.Operation{
		OperationID: "create-site",
		Method:      http.MethodPost,
		Path:        URLBuilder.Build("sites"),
		Tags:        []string{"Locations"},
		Summary:     "Add site",
	}, locationsControllers.CreateSite)
	huma.Register(*api, huma.Operation{
		OperationID: "create-camera",
		Method:      http.MethodPost,
		Path:        URLBuilder.Build("cameras"),
		Tags:        []string{"Locations"},
		Summary:     "Add camera",
	}, locationsControllers.AddCamera)
	huma.Register(*api, huma.Operation{
		OperationID: "delete-camera",
		Method:      http.MethodDelete,
		Path:        URLBuilder.Build("cameras/{id}"),
		Tags:        []string{"Locations"},
		Summary:     "Delete camera",
	}, locationsControllers.DeleteCamera)
	huma.Register(*api, huma.Operation{
		OperationID: "update-camera",
		Method:      http.MethodPut,
		Path:        URLBuilder.Build("cameras/{id}"),
		Tags:        []string{"Locations"},
		Summary:     "Update camera",
	}, locationsControllers.PutCamera)
	huma.Register(*api, huma.Operation{
		OperationID: "patch-camera",
		Method:      http.MethodPatch,
		Path:        URLBuilder.Build("cameras/{id}"),
		Tags:        []string{"Locations"},
		Summary:     "Patch camera",
	}, locationsControllers.PatchCamera)
	huma.Register(*api, huma.Operation{
		OperationID: "delete-site",
		Method:      http.MethodDelete,
		Path:        URLBuilder.Build("sites/{id}"),
		Tags:        []string{"Locations"},
		Summary:     "Delete site",
	}, locationsControllers.DeleteSite)
	huma.Register(*api, huma.Operation{
		OperationID: "put-site",
		Method:      http.MethodPut,
		Path:        URLBuilder.Build("sites/{id}"),
		Tags:        []string{"Locations"},
		Summary:     "Update site",
	}, locationsControllers.PutSite)

}
