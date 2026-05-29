package site

import (
	"fmt"

	api "github.com/lejeunel/go-image-annotator/adapters/api/server"
	web "github.com/lejeunel/go-image-annotator/adapters/web"
	a "github.com/lejeunel/go-image-annotator/app/annotator"
	"github.com/lejeunel/go-image-annotator/app/annotator/presenters"
	scr "github.com/lejeunel/go-image-annotator/app/annotator/scroller"

	"net/http"

	"github.com/lejeunel/go-image-annotator/config"
	"github.com/lejeunel/go-image-annotator/infra"
	i "github.com/lejeunel/go-image-annotator/infra/interactors"
)

type SiteConfig struct {
	APIDocsPath      string
	OpenAPISpecsPath string
}

func RegisterHandlers(mux *http.ServeMux, apiServer api.Server, webServer web.Server, cfg SiteConfig) {
	RegisterAPI(mux, apiServer, cfg.APIDocsPath, cfg.OpenAPISpecsPath)
	RegisterStaticFiles(mux)
	web.RegisterWebPages(mux, webServer)
}

func Serve(port int) {
	cfg := config.Parse()
	mux := http.NewServeMux()

	infra := infra.NewSQLiteInfra(cfg.DBPath, cfg.ArtefactDir)
	interactors := i.NewSQLiteInteractors(infra, cfg.DefaultPageSize, cfg.AllowedImageFormats)
	scroller := scr.New(infra.ScrollerRepo)
	annotator := a.NewAnnotator(scroller, &interactors.Image.Read,
		&interactors.Annotation.AddBox, &interactors.Annotation.UpdateBox, &interactors.Annotation.Delete,
		&interactors.Label.FetchAll, &interactors.Annotation.UpdateLabel, &interactors.Annotation.AddImageLabel,
		presenters.NewPresenter())
	RegisterHandlers(mux,
		*api.NewServer(interactors),
		*web.NewServer(interactors, annotator),
		SiteConfig{APIDocsPath: "/api/docs", OpenAPISpecsPath: "/api/openapi.yaml"})

	fmt.Println("serving on port:", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
}
