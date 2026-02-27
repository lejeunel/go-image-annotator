package labels

import (
	g "datahub/generic"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
)

func RegisterLabelsRoutes(api *huma.API, labelService *Service, URLBuilder g.APIURLBuilder) {
	c := &LabelsHTTPController{LabelsService: labelService}
	huma.Register(*api, huma.Operation{
		OperationID: "create-label",
		Method:      http.MethodPost,
		Path:        URLBuilder.Build("labels"),
		Tags:        []string{"Labels"},
		Summary:     "Add new label",
	}, c.Create)
	huma.Register(*api, huma.Operation{
		OperationID: "put-label",
		Method:      http.MethodPut,
		Path:        URLBuilder.Build("labels/{name}"),
		Tags:        []string{"Labels"},
		Summary:     "Update a label",
	}, c.Update)
	huma.Register(*api, huma.Operation{
		OperationID: "delete-label",
		Method:      http.MethodDelete,
		Path:        URLBuilder.Build("labels/{name}"),
		Tags:        []string{"Labels"},
		Summary:     "Delete label",
	}, c.Delete)
}
