package server

import (
	"fmt"
	"log/slog"
	"os"

	api "github.com/lejeunel/go-image-annotator/adapters/api/server"

	app "github.com/lejeunel/go-image-annotator/adapters/sqlite/app"
	"github.com/lejeunel/go-image-annotator/adapters/web"
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	a "github.com/lejeunel/go-image-annotator/app/annotator"
	ap "github.com/lejeunel/go-image-annotator/app/annotator/presenters"
	scr "github.com/lejeunel/go-image-annotator/app/annotator/scroller"
	as "github.com/lejeunel/go-image-annotator/assets"
	"github.com/lejeunel/go-image-annotator/config"
	ip "github.com/lejeunel/go-image-annotator/shared/identity_provider"
	sm "github.com/lejeunel/go-image-annotator/shared/session"

	"net/http"
)

func RegisterHandlers(mux *http.ServeMux, apiServer api.Server, webServer web.Server, apicfg config.APIConfig,
	pageBuilder b.PageBuilder) {
	api.RegisterAPIEndpoints(mux, apiServer, apicfg.APIPath)
	web.RegisterWebPages(mux, webServer, pageBuilder)
	web.RegisterAPIDocs(mux, apicfg.OpenAPISpecsPath, apicfg.APIDocsPath)
	as.RegisterAPISpecs(mux, apicfg.OpenAPISpecsPath)
	as.RegisterStaticFiles(mux)
}

func Make(apiPath string) http.Handler {
	app := app.NewSQLiteApp()
	mux := http.NewServeMux()

	pageBuilder := b.NewPageBuilder(apiPath)

	sessionManager := sm.NewSQLiteSessionManager(app.Infra.Db.DB, app.Infra.User, app.TokenGenerator)
	identityProvider := ip.NewGothIdentityHandler(sessionManager)
	ip.SetupForGoogle(ip.OAuthProviderConfig{Key: os.Getenv("GOIA_GOOGLE_CLIENT_ID"),
		Secret:      os.Getenv("GOIA_GOOGLE_CLIENT_SECRET"),
		CallbackURL: "http://localhost:3000/callback/google"})

	scroller := scr.New(app.Infra.Scroller)
	annotator := a.NewAnnotator(scroller, &app.Itrs.Image.Read,
		&app.Itrs.Annotation.AddBox, &app.Itrs.Annotation.UpdateBox, &app.Itrs.Annotation.Delete,
		&app.Itrs.Label.FetchAll, &app.Itrs.Annotation.UpdateLabel, &app.Itrs.Annotation.AddImageLabel,
		ap.New())
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	RegisterHandlers(mux,
		*api.NewServer(&app.Itrs, *logger),
		*web.NewServer(&app.Itrs, annotator, *pageBuilder, sessionManager, identityProvider),
		config.APIConfig{APIPath: fmt.Sprintf("/%v", apiPath),
			APIDocsPath:      fmt.Sprintf("/%v/docs", apiPath),
			OpenAPISpecsPath: fmt.Sprintf("/%v/openapi.yaml", apiPath)},
		*pageBuilder)

	return sessionManager.MiddleWare(mux)
}

func Serve(handler http.Handler, port int) {

	fmt.Println("serving on port:", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), handler)
}
