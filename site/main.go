package site

import (
	"fmt"
	"os"

	api "github.com/lejeunel/go-image-annotator/adapters/api/server"
	web "github.com/lejeunel/go-image-annotator/adapters/web"
	a "github.com/lejeunel/go-image-annotator/app/annotator"
	"github.com/lejeunel/go-image-annotator/app/annotator/presenters"
	scr "github.com/lejeunel/go-image-annotator/app/annotator/scroller"
	"github.com/lejeunel/go-image-annotator/shared/html"
	ip "github.com/lejeunel/go-image-annotator/shared/identity_provider"
	"github.com/lejeunel/go-image-annotator/shared/session"

	"net/http"

	"github.com/gorilla/sessions"
	"github.com/lejeunel/go-image-annotator/config"
	"github.com/lejeunel/go-image-annotator/infra"
	i "github.com/lejeunel/go-image-annotator/infra/interactors"
	"github.com/markbates/goth/gothic"
)

type SiteConfig struct {
	APIPath string
}

func Make(apiPath string) http.Handler {
	cfg := config.Parse()

	gothic.Store = sessions.NewCookieStore([]byte(os.Getenv("GOIA_SESSION_SECRET")))
	ip.SetupForGoogle(ip.OAuthProviderConfig{Key: os.Getenv("GOIA_GOOGLE_CLIENT_ID"),
		Secret:      os.Getenv("GOIA_GOOGLE_CLIENT_SECRET"),
		CallbackURL: "http://localhost:3000/callback/google"})
	mux := http.NewServeMux()

	infra := infra.NewSQLiteInfra(cfg.DBPath, cfg.ArtefactDir)

	sessionManager := session.NewSQLiteSessionManager(infra.Db.DB)
	identityHandler := ip.NewGothIdentityHandler(sessionManager)

	interactors := i.NewSQLiteInteractors(infra, cfg.DefaultPageSize, cfg.AllowedImageFormats)
	scroller := scr.New(infra.ScrollerRepo)
	annotator := a.NewAnnotator(scroller, &interactors.Image.Read,
		&interactors.Annotation.AddBox, &interactors.Annotation.UpdateBox, &interactors.Annotation.Delete,
		&interactors.Label.FetchAll, &interactors.Annotation.UpdateLabel, &interactors.Annotation.AddImageLabel,
		presenters.NewPresenter())

	htmlPageBuilder := html.NewPageBuilder("api")
	api.HandlerFromMuxWithBaseURL(api.NewServer(interactors), mux, fmt.Sprintf("/%v", apiPath))
	RegisterAPISpecs(mux, apiPath)
	RegisterStaticFiles(mux)
	web.RegisterWebPages(mux,
		*web.NewServer(interactors, annotator,
			*htmlPageBuilder, sessionManager, identityHandler), *htmlPageBuilder)

	return sessionManager.MiddleWare(mux)
}

func Serve(port int, handler http.Handler) {

	fmt.Println("serving on port:", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), handler)
}
