package server

import (
	"fmt"
	"log/slog"
	"os"

	api "github.com/lejeunel/go-image-annotator/adapters/api/server"
	userDashboard "github.com/lejeunel/go-image-annotator/adapters/web/user-dashboard"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	rt "github.com/lejeunel/go-image-annotator/routes"

	"github.com/lejeunel/go-image-annotator/adapters/web"
	adm "github.com/lejeunel/go-image-annotator/adapters/web/admin"
	admgrp "github.com/lejeunel/go-image-annotator/adapters/web/admin/group"
	admusr "github.com/lejeunel/go-image-annotator/adapters/web/admin/user"
	an "github.com/lejeunel/go-image-annotator/adapters/web/annotator"
	wauth "github.com/lejeunel/go-image-annotator/adapters/web/auth"
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	clc "github.com/lejeunel/go-image-annotator/adapters/web/collection"
	im "github.com/lejeunel/go-image-annotator/adapters/web/image"
	lbl "github.com/lejeunel/go-image-annotator/adapters/web/label"
	a "github.com/lejeunel/go-image-annotator/app"
	"github.com/lejeunel/go-image-annotator/app/sqlite"
	"github.com/lejeunel/go-image-annotator/config"
	g "github.com/lejeunel/go-image-annotator/globals"

	"net/http"

	"github.com/go-chi/chi/v5"
)

func Make(auth auth.Authorizer, url string, port int) http.Handler {
	cfg := config.Parse()

	currentVersion := g.Info{Version: g.Version, Commit: g.Commit, Date: g.Date}
	basePageBuilder := b.NewBasePageBuilder()
	pageBuilder := b.NewPageBuilder(basePageBuilder, currentVersion)

	app := sqlite.NewSQLiteApp(cfg, auth)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	a.MaybeCreateInitialAdmin(app.Itrs.User.Create, cfg.InitialAdminEmail, cfg.InitialAdminPassword)

	baseURL := fmt.Sprintf("%v:%v", url, port)

	router := chi.NewRouter()
	webAuth := Chain(
		app.SessionManager.LoadAndSave,
		app.SessionManager.AuthCookiesMiddleWare,
		WebRequireLogin,
	)

	RouteWebPages(
		router,
		*web.NewServer(&app.Itrs, pageBuilder, app.SessionManager, cfg.DefaultPageSize),
		HomePageHandlerFunc(pageBuilder),
		webAuth,
	)
	udb := userDashboard.New(pageBuilder, app.Itrs.User.RenewToken, app.Itrs.User.ChangePassword)
	udb.Route(router, webAuth)

	RouteAPI(router, *api.NewServer(&app.Itrs, *logger),
		app.SessionManager.LoadAndSave, app.SessionManager.AuthBearerMiddleWare, app.SessionManager.AuthCookiesMiddleWare, ApiRequireLogin)
	RouteAPIDocs(router, APIDocsHandlerFunc(rt.APISpecs, pageBuilder), webAuth)
	RouteAPISpecs(router)
	RouteStaticFiles(router)

	annotatorServer := an.NewServer(app.Annotator, pageBuilder, app.SessionManager)
	annotatorServer.Route(router, webAuth)

	collectionServer := clc.New(pageBuilder, cfg.DefaultPageSize,
		app.Itrs.Collection.Create, app.Itrs.Collection.List, app.Itrs.Collection.Update,
		app.Itrs.Collection.Delete, app.Itrs.Collection.Find)
	collectionServer.Route(router, webAuth)

	imagesServer := im.New(pageBuilder, cfg.DefaultPageSize, app.Itrs.Image.List, app.Itrs.Image.Delete, app.Itrs.Image.Find)
	imagesServer.Route(router, webAuth)

	adminPageBuilder := adm.NewPageBuilder(pageBuilder)
	adminUserServer := admusr.New(adminPageBuilder, app.Itrs.User, app.Itrs.Group, app.Itrs.Role, cfg.DefaultPageSize)
	adminUserServer.Route(router, webAuth)
	adminGroupServer := admgrp.New(adminPageBuilder, app.Itrs.Group)
	adminGroupServer.Route(router, webAuth)

	labelServer := lbl.New(pageBuilder, cfg.DefaultPageSize,
		app.Itrs.Label.Create, app.Itrs.Label.List, app.Itrs.Label.Update,
		app.Itrs.Label.Delete, app.Itrs.Label.Find)
	labelServer.Route(router, webAuth)

	authServer := wauth.New(
		baseURL,
		basePageBuilder,
		*logger,
		app.SessionManager,
		app.Itrs.User.RequestForgottenPassword,
		app.Itrs.User.ResetForgottenPassword)
	authServer.Route(router,
		app.SessionManager.LoadAndSave)

	return router
}

func Serve(handler http.Handler, port int) {

	fmt.Println("serving on port:", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), handler)
}
