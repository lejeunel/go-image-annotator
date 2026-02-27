package app

import (
	id "datahub/app/identity"
	in "datahub/app/ingester"
	pro "datahub/domain/annotation_profiles"
	clc "datahub/domain/collections"
	im "datahub/domain/images"
	lbl "datahub/domain/labels"
	loc "datahub/domain/locations"
	g "datahub/generic"
	"fmt"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"net/http"
)

func RegisterAPIRoutes(app *App, mux *http.ServeMux, version string) {
	apiCfg := huma.DefaultConfig("Datahub API", "1.0.0")
	APIURLBuilder := g.NewAPIURLBuilder(version)
	apiCfg.DocsPath = APIURLBuilder.Build("{$}")
	api := humago.New(mux, apiCfg)
	AddAPIRoutes(&api, mux, app, APIURLBuilder)
}

func AddAPIRoutes(api *huma.API, mux *http.ServeMux, app *App, URLBuilder *g.APIURLBuilder) {

	identityController := &id.IdentityHTTPController{Authorizer: app.Authorizer}
	rawImageURLBuilder := im.NewRawImageURLBuilder(URLBuilder, "raw-image")
	imagesController := im.NewImagesHTTPController(app.Images, app.Collections,
		app.Locations, rawImageURLBuilder)

	url := fmt.Sprintf("%v/{id}", URLBuilder.Build("raw-image"))
	mux.HandleFunc(url, imagesController.GetRaw)
	im.RegisterImagesRoutes(api, imagesController, *URLBuilder)
	im.RegisterBaseImageRoutes(api, imagesController, URLBuilder)
	clc.RegisterCollectionsRoutes(api, app.Collections, *URLBuilder)
	pro.RegisterProfilesRoutes(api, app.Profiles, *URLBuilder)
	loc.RegisterLocationsRoutes(api, app.Locations, *URLBuilder)
	lbl.RegisterLabelsRoutes(api, app.Labels, *URLBuilder)
	in.RegisterIngestionRoutes(api, app.Ingestion, *URLBuilder)

	huma.Register(*api, huma.Operation{
		OperationID: "get-identity",
		Method:      http.MethodGet,
		Path:        URLBuilder.Build("whoami"),
		Tags:        []string{"Identity"},
		Summary:     "Get Identity",
	}, identityController.Get)

}
