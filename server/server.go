package server

import (
	"fmt"
	"log/slog"
	"os"

	api "github.com/lejeunel/go-image-annotator/adapters/api/server"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	rt "github.com/lejeunel/go-image-annotator/server/routes"

	"github.com/lejeunel/go-image-annotator/adapters/web"
	ap "github.com/lejeunel/go-image-annotator/adapters/web/annotator/presenters"
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	a "github.com/lejeunel/go-image-annotator/app"
	"github.com/lejeunel/go-image-annotator/app/sqlite"
	"github.com/lejeunel/go-image-annotator/config"
	ip "github.com/lejeunel/go-image-annotator/shared/identity_provider"

	"github.com/go-chi/chi/v5"
	"net/http"
)

func Make(auth auth.Authorizer) http.Handler {
	cfg := config.Parse()

	basePageBuilder := b.NewBasePageBuilder()
	basePageBuilder.AddScripts(b.BaseLibs()...)
	pageBuilder := b.NewPageBuilder(basePageBuilder, cfg.APIPath, cfg.RepoURL, cfg.DocsURL)
	loginPageBuilder := b.NewLoginPageBuilder(basePageBuilder)
	forgotPasswordPageBuilder := b.NewForgotPasswordBuilder(basePageBuilder)

	app := sqlite.NewSQLiteApp(cfg, auth)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	a.MaybeCreateInitialAdmin(app.Itrs.User.Create, cfg.InitialAdminEmail, cfg.InitialAdminPassword)

	router := chi.NewRouter()

	colorizer := ap.NewCyclicColorizer(ap.Palette)

	ip.SetupForGoogle(ip.OAuthProviderConfig{Key: os.Getenv("GOIA_GOOGLE_CLIENT_ID"),
		Secret:      os.Getenv("GOIA_GOOGLE_CLIENT_SECRET"),
		CallbackURL: "http://localhost:3000/auth/callback/google"})
	loginPageBuilder.AddOAuthProvider("google", "/auth/login/google")

	rt.RouteWebPages(
		router,
		*web.NewServer(&app.Itrs, app.Annotator,
			*pageBuilder, ap.NewAnnotationPagePresenter(colorizer),
			ap.NewAnnotoriousPresenter(colorizer),
			app.SessionManager, cfg.DefaultPageSize),
		HomePageHandlerFunc(*pageBuilder),
		app.SessionManager.LoadAndSave, app.SessionManager.AuthCookiesMiddleWare, WebRequireLogin,
	)
	rt.RouteAPI(router, *api.NewServer(&app.Itrs, *logger),
		app.SessionManager.LoadAndSave, app.SessionManager.AuthBearerMiddleWare, app.SessionManager.AuthCookiesMiddleWare, ApiRequireLogin)
	rt.RouteAPIDocs(router, APIDocsHandlerFunc(rt.APISpecs, *pageBuilder),
		app.SessionManager.LoadAndSave, app.SessionManager.AuthCookiesMiddleWare, WebRequireLogin,
	)
	rt.RouteAPISpecs(router)
	rt.RouteStaticFiles(router)
	rt.RouteAuth(router, app.AuthHandler, LoginPageHandlerFunc(*loginPageBuilder),
		ForgotPasswordHandlerFunc(*forgotPasswordPageBuilder),
		app.SessionManager.LoadAndSave)

	return router
}

func Serve(handler http.Handler, port int) {

	fmt.Println("serving on port:", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), handler)
}
