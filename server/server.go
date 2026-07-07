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
	a "github.com/lejeunel/go-image-annotator/app"
	"github.com/lejeunel/go-image-annotator/app/sqlite"
	as "github.com/lejeunel/go-image-annotator/assets"
	"github.com/lejeunel/go-image-annotator/config"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	ip "github.com/lejeunel/go-image-annotator/shared/identity_provider"

	"net/http"
)

func LoginPageHandlerFunc(builder b.LoginPageBuilder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		builder.Render(w)
	}
}
func WebRequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := u.IdentityFromContext(r.Context())
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func ApiRequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := u.IdentityFromContext(r.Context())
		if user == nil {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func MakeWebPagesMux(webServer web.Server, apicfg config.APIConfig,
	pageBuilder b.PageBuilder) *http.ServeMux {
	mux := http.NewServeMux()
	web.RegisterWebPages(mux, webServer, pageBuilder)
	web.RegisterAPIDocs(mux, apicfg.OpenAPISpecsPath, apicfg.APIDocsPath, pageBuilder)
	as.RegisterAPISpecs(mux, apicfg.OpenAPISpecsPath)
	return mux

}

func MakeAPIMux(apiServer api.Server) *http.ServeMux {
	mux := http.NewServeMux()
	api.HandlerFromMux(&apiServer, mux)
	return mux
}

func Make(auth auth.Auth) http.Handler {
	cfg := config.Parse()

	pageBuilder := b.NewPageBuilder(cfg.APIPath, cfg.RepoURL, cfg.DocsURL)
	loginPageBuilder := b.NewLoginPageBuilder()

	app := sqlite.NewSQLiteApp(cfg, auth)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	a.MaybeCreateInitialAdmin(app.Itrs.User.Create, cfg.InitialAdminEmail, cfg.InitialAdminPassword)

	ip.SetupForGoogle(ip.OAuthProviderConfig{Key: os.Getenv("GOIA_GOOGLE_CLIENT_ID"),
		Secret:      os.Getenv("GOIA_GOOGLE_CLIENT_SECRET"),
		CallbackURL: "http://localhost:3000/auth/callback/google"})
	loginPageBuilder.AddOAuthProvider("google", "/auth/login/google")
	apiCfg := config.APIConfig{APIPath: fmt.Sprintf("/%v", cfg.APIPath),
		APIDocsPath:      fmt.Sprintf("/%v/docs", cfg.APIPath),
		OpenAPISpecsPath: fmt.Sprintf("/%v/openapi.yaml", cfg.APIPath)}
	colorizer := ap.NewCyclicColorizer(ap.Palette)
	webPagesMux := MakeWebPagesMux(
		*web.NewServer(&app.Itrs, app.Annotator,
			*pageBuilder, ap.NewAnnotationPagePresenter(colorizer),
			ap.NewAnnotoriousPresenter(colorizer),
			app.SessionManager, cfg.DefaultPageSize),
		apiCfg,
		*pageBuilder,
	)
	apiMux := MakeAPIMux(*api.NewServer(&app.Itrs, *logger))
	oauthMux := web.MakeOAuthMux(app.OAuthHandler)

	rootMux := http.NewServeMux()

	rootMux.Handle("/auth/", http.StripPrefix("/auth", app.SessionManager.LoadAndSave(oauthMux)))
	rootMux.Handle("/login/", LoginPageHandlerFunc(*loginPageBuilder))

	rootMux.Handle("/api/", http.StripPrefix("/api", app.SessionManager.ApiMiddleWare(ApiRequireLogin(apiMux))))
	as.RegisterStaticFiles(rootMux)
	rootMux.Handle("/", app.SessionManager.WebPagesMiddleWare(WebRequireLogin(webPagesMux)))

	return rootMux
}

func Serve(handler http.Handler, port int) {

	fmt.Println("serving on port:", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), handler)
}
