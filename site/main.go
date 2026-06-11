package site

import (
	"fmt"
	"log/slog"
	"os"

	api "github.com/lejeunel/go-image-annotator/adapters/api/server"
	web "github.com/lejeunel/go-image-annotator/adapters/web"
	a "github.com/lejeunel/go-image-annotator/app/annotator"
	"github.com/lejeunel/go-image-annotator/app/annotator/presenters"
	scr "github.com/lejeunel/go-image-annotator/app/annotator/scroller"
	tok "github.com/lejeunel/go-image-annotator/app/token"
	"github.com/lejeunel/go-image-annotator/shared/html"
	ip "github.com/lejeunel/go-image-annotator/shared/identity_provider"
	sm "github.com/lejeunel/go-image-annotator/shared/session"

	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/sqlite"
	i "github.com/lejeunel/go-image-annotator/adapters/sqlite/interactors"
	"github.com/lejeunel/go-image-annotator/config"
)

type SiteConfig struct {
	APIPath          string
	APIDocsPath      string
	OpenAPISpecsPath string
}

func RegisterHandlers(mux *http.ServeMux, apiServer api.Server, webServer web.Server, cfg SiteConfig,
	pageBuilder html.PageBuilder) {
	RegisterAPIDocs(mux, cfg.OpenAPISpecsPath, cfg.APIDocsPath)
	RegisterAPISpecs(mux, cfg.OpenAPISpecsPath)
	RegisterAPIEndpoints(mux, apiServer, cfg.APIPath)
	RegisterStaticFiles(mux)
	web.RegisterWebPages(mux, webServer, pageBuilder)
}

func Make(apiPath string) http.Handler {
	cfg := config.Parse()
	mux := http.NewServeMux()

	infra := infra.NewSQLiteInfra(cfg.DBPath, cfg.ArtefactDir)
	tokenGenerator := tok.NewTokenGenerator(32)
	interactors := i.NewSQLiteInteractors(infra, cfg.DefaultPageSize, cfg.AllowedImageFormats, tokenGenerator)
	pageBuilder := html.NewPageBuilder(apiPath)

	sessionManager := sm.NewSQLiteSessionManager(infra.Db.DB, infra.User, tokenGenerator)
	identityProvider := ip.NewGothIdentityHandler(sessionManager)
	ip.SetupForGoogle(ip.OAuthProviderConfig{Key: os.Getenv("GOIA_GOOGLE_CLIENT_ID"),
		Secret:      os.Getenv("GOIA_GOOGLE_CLIENT_SECRET"),
		CallbackURL: "http://localhost:3000/callback/google"})

	scroller := scr.New(infra.Scroller)
	annotator := a.NewAnnotator(scroller, &interactors.Image.Read,
		&interactors.Annotation.AddBox, &interactors.Annotation.UpdateBox, &interactors.Annotation.Delete,
		&interactors.Label.FetchAll, &interactors.Annotation.UpdateLabel, &interactors.Annotation.AddImageLabel,
		presenters.NewPresenter())
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	RegisterHandlers(mux,
		*api.NewServer(interactors, *logger),
		*web.NewServer(interactors, annotator, *pageBuilder, sessionManager, identityProvider),
		SiteConfig{APIPath: fmt.Sprintf("/%v", apiPath),
			APIDocsPath:      fmt.Sprintf("/%v/docs", apiPath),
			OpenAPISpecsPath: fmt.Sprintf("/%v/openapi.yaml", apiPath)},
		*pageBuilder)

	return sessionManager.MiddleWare(mux)
}

func Serve(handler http.Handler, port int) {

	fmt.Println("serving on port:", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), handler)
}
