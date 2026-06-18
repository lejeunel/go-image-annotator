package server

import (
	"fmt"
	"log/slog"
	"os"

	api "github.com/lejeunel/go-image-annotator/adapters/api/server"
	"github.com/lejeunel/go-image-annotator/modules/auth"

	"github.com/lejeunel/go-image-annotator/adapters/web"
	ap "github.com/lejeunel/go-image-annotator/adapters/web/annotator/presenters"
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	"github.com/lejeunel/go-image-annotator/app/sqlite"
	as "github.com/lejeunel/go-image-annotator/assets"
	"github.com/lejeunel/go-image-annotator/config"
	ip "github.com/lejeunel/go-image-annotator/shared/identity_provider"

	"net/http"
)

func RegisterHandlers(mux *http.ServeMux, apiServer api.Server, webServer web.Server, apicfg config.APIConfig,
	pageBuilder b.PageBuilder) {
	api.RegisterAPIEndpoints(mux, apiServer, apicfg.APIPath)
	web.RegisterWebPages(mux, webServer, pageBuilder)
	web.RegisterAPIDocs(mux, apicfg.OpenAPISpecsPath, apicfg.APIDocsPath, pageBuilder)
	as.RegisterAPISpecs(mux, apicfg.OpenAPISpecsPath)
	as.RegisterStaticFiles(mux)
}

func Make(auth auth.Auth) http.Handler {
	cfg := config.Parse()
	mux := http.NewServeMux()

	pageBuilder := b.NewPageBuilder(cfg.APIPath, cfg.RepoURL, cfg.DocsURL)

	app := sqlite.NewSQLiteApp(cfg, auth)
	ip.SetupForGoogle(ip.OAuthProviderConfig{Key: os.Getenv("GOIA_GOOGLE_CLIENT_ID"),
		Secret:      os.Getenv("GOIA_GOOGLE_CLIENT_SECRET"),
		CallbackURL: "http://localhost:3000/callback/google"})

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	colorizer := ap.NewCyclicColorizer(ap.Palette)
	RegisterHandlers(mux,
		*api.NewServer(&app.Itrs, *logger),
		*web.NewServer(&app.Itrs, app.Annotator,
			*pageBuilder, ap.NewAnnotationPagePresenter(colorizer),
			ap.NewAnnotoriousPresenter(colorizer),
			app.SessionManager, app.OAuthHandler),
		config.APIConfig{APIPath: fmt.Sprintf("/%v", cfg.APIPath),
			APIDocsPath:      fmt.Sprintf("/%v/docs", cfg.APIPath),
			OpenAPISpecsPath: fmt.Sprintf("/%v/openapi.yaml", cfg.APIPath)},
		*pageBuilder)

	return app.SessionManager.MiddleWare(mux)
}

func Serve(handler http.Handler, port int) {

	fmt.Println("serving on port:", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), handler)
}
