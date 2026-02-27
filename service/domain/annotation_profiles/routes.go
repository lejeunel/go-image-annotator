package annotation_profile

import (
	g "datahub/generic"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
)

func RegisterProfilesRoutes(api *huma.API, s *Service, URLBuilder g.APIURLBuilder) {
	profilesController := &ProfileHTTPController{ProfileService: s}
	huma.Register(*api, huma.Operation{
		OperationID: "get-profile",
		Method:      http.MethodGet,
		Path:        URLBuilder.Build("profiles/{name}"),
		Tags:        []string{"Annotation Profiles"},
		Summary:     "Get profile",
		Description: "Get profile",
	}, profilesController.Get)
	huma.Register(*api, huma.Operation{
		OperationID: "list-profile",
		Method:      http.MethodGet,
		Path:        URLBuilder.Build("profiles"),
		Tags:        []string{"Annotation Profiles"},
		Summary:     "List profiles",
		Description: "List profiles",
	}, profilesController.List)
	huma.Register(*api, huma.Operation{
		OperationID: "create-profile",
		Method:      http.MethodPost,
		Path:        URLBuilder.Build("profiles"),
		Tags:        []string{"Annotation Profiles"},
		Summary:     "Create a profile",
	}, profilesController.Create)
	huma.Register(*api, huma.Operation{
		OperationID: "delete-profile",
		Method:      http.MethodDelete,
		Path:        URLBuilder.Build("profiles/{name}"),
		Tags:        []string{"Annotation Profiles"},
		Summary:     "Delete a profile",
	}, profilesController.Delete)
	huma.Register(*api, huma.Operation{
		OperationID: "put-profile",
		Method:      http.MethodPut,
		Path:        URLBuilder.Build("profiles/{name}"),
		Tags:        []string{"Annotation Profiles"},
		Summary:     "Update a profile",
	}, profilesController.Put)
	huma.Register(*api, huma.Operation{
		OperationID: "add-label-to-profile",
		Method:      http.MethodPost,
		Path:        URLBuilder.Build("profiles/{profile}/label/{label}"),
		Tags:        []string{"Annotation Profiles"},
		Summary:     "Add label to profile",
	}, profilesController.AddLabel)
	huma.Register(*api, huma.Operation{
		OperationID: "remove-label-from-profile",
		Method:      http.MethodDelete,
		Path:        URLBuilder.Build("profiles/{profile}/label/{label}"),
		Tags:        []string{"Annotation Profiles"},
		Summary:     "Remove label from profile",
	}, profilesController.RemoveLabel)

}
