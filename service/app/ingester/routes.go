package ingestion

import (
	g "datahub/generic"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
)

func RegisterIngestionRoutes(api *huma.API, ingestionService *Service, URLBuilder g.APIURLBuilder) {
	ingestionController := &IngestionHTTPController{IngestionService: ingestionService}
	huma.Register(*api, huma.Operation{
		OperationID:  "create-image",
		Method:       http.MethodPost,
		Path:         URLBuilder.Build("collections/{name}/images"),
		Tags:         []string{"Images"},
		Summary:      "Ingest image",
		MaxBodyBytes: 10 * 1024 * 1024, //10MB
	}, ingestionController.Ingest)
}
