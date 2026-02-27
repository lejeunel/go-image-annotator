package site

import (
	"context"
	a "datahub/app"
	an "datahub/app/annotator"
	au "datahub/app/authorizer"
	m "datahub/app/migrations"
	c "datahub/config"
	pro "datahub/domain/annotation_profiles"
	clc "datahub/domain/collections"
	im "datahub/domain/images"
	lbl "datahub/domain/labels"
	loc "datahub/domain/locations"
	g "datahub/generic"
	"embed"
	"fmt"
	clk "github.com/jonboulle/clockwork"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//go:embed static
var staticFiles embed.FS

//go:embed docs
var docsFiles embed.FS

// func MakeDocsMux() (*http.ServeMux, error) {

// 	docsFileSystem := fs.FS(docsFiles)
// 	docsSubFileSystem, err := fs.Sub(docsFileSystem, "docs/public")
// 	if err != nil {
// 		log.Fatal("Error building docs filesystem")
// 		return nil, err
// 	}
// 	httpDocsFileServer := http.FileServer(http.FS(docsSubFileSystem))
// 	docsMux := http.NewServeMux()
// 	docsMux.Handle("/", httpDocsFileServer)

// 	return docsMux, nil

// }

func MakeStaticMux() (*http.ServeMux, error) {
	staticFileSystem := fs.FS(staticFiles)
	staticSubFileSystem, err := fs.Sub(staticFileSystem, "static")
	if err != nil {
		log.Fatal("Error building static filesystem")
		return nil, err
	}
	httpStaticFileServer := http.FileServer(http.FS(staticSubFileSystem))
	staticMux := http.NewServeMux()
	staticMux.Handle("/", httpStaticFileServer)

	return staticMux, nil

}

func RegisterViewsHandlers(app *a.App, mux *http.ServeMux) {
	genericViewer := &g.GenericViewer{SignOutURL: app.Config.SignOutURL,
		IdentityProvider: app.Authorizer.IdentityProvider}

	collectionListViewer := clc.CollectionsListViewer{Viewer: genericViewer,
		CollectionService: app.Collections, PageSize: 10, PaginationWidgetSize: 10}
	labelViewer := lbl.LabelListViewer{Viewer: genericViewer, LabelsService: app.Labels,
		PageSize:             10,
		PaginationWidgetSize: 10}
	siteListViewer := loc.SitesListViewer{Viewer: genericViewer,
		LocationService: app.Locations, PageSize: 10,
		PaginationWidgetSize: 10,
	}
	siteViewer := loc.CamerasListViewer{Viewer: genericViewer,
		LocationService: app.Locations}
	profileListViewer := pro.AnnotationProfilesListViewer{Viewer: genericViewer,
		Service: app.Profiles, PageSize: 10, PaginationWidgetSize: 10,
	}

	imageAnnotatorModel := an.NewAnnotator(app.Labels,
		app.Images, app.Collections, app.Locations,
		app.Authorizer, app.Logger,
		app.Config.TargetImageWidth)
	imageAnnotatorViewer := an.NewAnnotatorViewer(genericViewer)

	imageAnnotatorController := an.NewAnnotatorController(imageAnnotatorModel,
		imageAnnotatorViewer, app.Logger)
	imageAnnotatorController.RegisterEndPoints(mux)

	collectionHandler := im.NewImageListOfCollectionHandler(genericViewer, app.Images, app.Collections, app.Images.MaxPageSize)
	cameraHandler := im.NewImageListOfCameraHandler(genericViewer, app.Images, app.Locations, app.Images.MaxPageSize)
	labelHandler := im.NewImageListOfLabelHandler(genericViewer, app.Images, app.Labels, app.Images.MaxPageSize)
	mux.Handle("/{$}", collectionListViewer.Handler())
	mux.Handle("/collections", collectionListViewer.Handler())
	mux.Handle("/collection/{id}", collectionHandler)
	mux.Handle("/labels", labelViewer.Handler())
	mux.Handle("/sites", siteListViewer.Handler())
	mux.Handle("/site/{id}", siteViewer.Handler())
	mux.Handle("/camera/{id}", cameraHandler)
	mux.Handle("/label/{id}", labelHandler)
	mux.Handle("/profiles", profileListViewer.Handler())
}

func PrependMiddlewaresToServeMux(r *http.ServeMux, middlewares ...func(next http.Handler) http.Handler) http.Handler {
	var s http.Handler
	s = r

	for _, mw := range middlewares {
		s = mw(s)
	}

	return s
}

func Serve(port int, migrate bool) error {
	cfg := c.NewConfig()

	app, db, migrationProvider, ctx := a.NewApp(cfg, clk.NewRealClock(), 1)

	if migrate == true {
		app.Logger.Info("Applying migrations...")
		if err := m.ApplyMigrations(ctx, migrationProvider, "up", app.Logger); err != nil {
			app.Logger.Error("Failed to run migrations on db %v", "db", db, "error", err)
			panic(err)
		}
		app.Logger.Info("Done applying migrations.")
	}

	// Create a new router & API.
	mainMux := http.NewServeMux()

	a.RegisterAPIRoutes(&app, mainMux, "v1")

	// docsMux, err := MakeDocsMux()
	// if err != nil {
	// 	return err
	// }

	staticMux, err := MakeStaticMux()
	if err != nil {
		return err
	}

	// mainMux.Handle("/docs/", http.StripPrefix("/docs", docsMux))
	mainMux.Handle("/static/", http.StripPrefix("/static", staticMux))
	RegisterViewsHandlers(&app, mainMux)

	var server *http.Server
	if cfg.Mode != "test" {
		app.Logger.Info("Applying authentication middleware")
		authMux := PrependMiddlewaresToServeMux(mainMux, au.AuthentikHeadersMiddleware)
		server = &http.Server{Addr: fmt.Sprintf(":%v", port), Handler: authMux}
	} else {
		server = &http.Server{Addr: fmt.Sprintf(":%v", port), Handler: mainMux}
	}

	// Start the server!
	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	app.Logger.Info("Listening", "port", port)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	log.Println("server stopping...")
	defer cancel()

	log.Fatal(server.Shutdown(ctx))

	return nil

}
