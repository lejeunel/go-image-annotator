package server

import (
	"fmt"
	"log/slog"
	"os"

	api "github.com/lejeunel/go-image-annotator/adapters/api/server"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	rt "github.com/lejeunel/go-image-annotator/routes"

	"github.com/lejeunel/go-image-annotator/adapters/web"
	ap "github.com/lejeunel/go-image-annotator/adapters/web/annotator/presenters"
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	a "github.com/lejeunel/go-image-annotator/app"
	"github.com/lejeunel/go-image-annotator/app/sqlite"
	"github.com/lejeunel/go-image-annotator/config"
	ip "github.com/lejeunel/go-image-annotator/shared/identity_provider"

	"net/http"

	"github.com/go-chi/chi/v5"
)

func Make(auth auth.Authorizer, url string, port int) http.Handler {
	cfg := config.Parse()

	basePageBuilder := b.NewBasePageBuilder()
	pageBuilder := b.NewPageBuilder(basePageBuilder, rt.APIRoot, cfg.RepoURL, cfg.DocsURL)
	loginPageBuilder := b.NewLoginPageBuilder(basePageBuilder)
	forgotPasswordPageBuilder := b.NewForgotPasswordBuilder(basePageBuilder)

	app := sqlite.NewSQLiteApp(cfg, auth)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	a.MaybeCreateInitialAdmin(app.Itrs.User.Create, cfg.InitialAdminEmail, cfg.InitialAdminPassword)

	baseURL := fmt.Sprintf("%v:%v", url, port)

	ip.SetupForGoogle(ip.OAuthProviderConfig{Key: os.Getenv("GOIA_GOOGLE_CLIENT_ID"),
		Secret:      os.Getenv("GOIA_GOOGLE_CLIENT_SECRET"),
		CallbackURL: rt.MakeOAuthCallbackURL(baseURL, "google")})
	loginPageBuilder.AddOAuthProvider("google", rt.MakeOAuthLoginURL("google"))

	router := chi.NewRouter()
	colorizer := ap.NewCyclicColorizer(ap.Palette)
	RouteWebPages(
		router,
		*web.NewServer(&app.Itrs, app.Annotator,
			*pageBuilder, ap.NewAnnotationPagePresenter(colorizer),
			ap.NewAnnotoriousPresenter(colorizer),
			app.SessionManager, cfg.DefaultPageSize),
		HomePageHandlerFunc(*pageBuilder),
		app.SessionManager.LoadAndSave, app.SessionManager.AuthCookiesMiddleWare, WebRequireLogin,
	)
	RouteAPI(router, *api.NewServer(&app.Itrs, *logger),
		app.SessionManager.LoadAndSave, app.SessionManager.AuthBearerMiddleWare, app.SessionManager.AuthCookiesMiddleWare, ApiRequireLogin)
	RouteAPIDocs(router, APIDocsHandlerFunc(rt.APISpecs, *pageBuilder),
		app.SessionManager.LoadAndSave, app.SessionManager.AuthCookiesMiddleWare, WebRequireLogin,
	)
	RouteAPISpecs(router)
	RouteStaticFiles(router)
	RouteAuth(router, app.AuthHandler, LoginPageHandlerFunc(*loginPageBuilder),
		ForgotPasswordHandlerFunc(*forgotPasswordPageBuilder),
		app.SessionManager.LoadAndSave)

	return router
}

func Serve(handler http.Handler, port int) {

	fmt.Println("serving on port:", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), handler)
}
