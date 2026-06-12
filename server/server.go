package server

import (
	"fmt"
	"log/slog"
	"os"

	api "github.com/lejeunel/go-image-annotator/adapters/api/server"
	db "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos"
	"github.com/lejeunel/go-image-annotator/adapters/web"
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	a "github.com/lejeunel/go-image-annotator/app/annotator"
	ap "github.com/lejeunel/go-image-annotator/app/annotator/presenters"
	scr "github.com/lejeunel/go-image-annotator/app/annotator/scroller"
	fs "github.com/lejeunel/go-image-annotator/app/file-store"
	tok "github.com/lejeunel/go-image-annotator/app/token"
	as "github.com/lejeunel/go-image-annotator/assets"
	"github.com/lejeunel/go-image-annotator/config"
	ip "github.com/lejeunel/go-image-annotator/shared/identity_provider"
	sm "github.com/lejeunel/go-image-annotator/shared/session"

	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/sqlite/infra"
	i "github.com/lejeunel/go-image-annotator/adapters/sqlite/interactors"
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
	cfg := config.Parse()
	mux := http.NewServeMux()

	infra := infra.NewSQLiteInfra(db.NewSQLiteDB(cfg.DBPath),
		fs.NewFileStore(cfg.ArtefactDir))
	tokenGenerator := tok.NewTokenGenerator(32)
	interactors := i.NewSQLiteInteractors(infra, cfg.DefaultPageSize, cfg.AllowedImageFormats, tokenGenerator)
	pageBuilder := b.NewPageBuilder(apiPath)

	sessionManager := sm.NewSQLiteSessionManager(infra.Db.DB, infra.User, tokenGenerator)
	identityProvider := ip.NewGothIdentityHandler(sessionManager)
	ip.SetupForGoogle(ip.OAuthProviderConfig{Key: os.Getenv("GOIA_GOOGLE_CLIENT_ID"),
		Secret:      os.Getenv("GOIA_GOOGLE_CLIENT_SECRET"),
		CallbackURL: "http://localhost:3000/callback/google"})

	scroller := scr.New(infra.Scroller)
	annotator := a.NewAnnotator(scroller, &interactors.Image.Read,
		&interactors.Annotation.AddBox, &interactors.Annotation.UpdateBox, &interactors.Annotation.Delete,
		&interactors.Label.FetchAll, &interactors.Annotation.UpdateLabel, &interactors.Annotation.AddImageLabel,
		ap.New())
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	RegisterHandlers(mux,
		*api.NewServer(interactors, *logger),
		*web.NewServer(interactors, annotator, *pageBuilder, sessionManager, identityProvider),
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
