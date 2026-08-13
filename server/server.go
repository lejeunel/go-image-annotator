package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	api "github.com/lejeunel/go-image-annotator/adapters/api/server"
	userDashboard "github.com/lejeunel/go-image-annotator/adapters/web/dashboard"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	rt "github.com/lejeunel/go-image-annotator/routes"

	adm "github.com/lejeunel/go-image-annotator/adapters/web/admin"
	admgrp "github.com/lejeunel/go-image-annotator/adapters/web/admin/group"
	admpl "github.com/lejeunel/go-image-annotator/adapters/web/admin/policy"
	admrl "github.com/lejeunel/go-image-annotator/adapters/web/admin/role"
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

	"github.com/go-chi/chi/v5"
)

func Make(port int) http.Handler {
	cfg := config.Parse()
	defaultAuth := auth.NewDefault()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	app := sqlite.NewApp(cfg, &defaultAuth, *logger)

	currentVersion := g.Info{Version: g.Version, Date: g.Date}
	basePageBuilder := b.NewBasePageBuilder()
	pageBuilder := b.NewPageBuilder(basePageBuilder, currentVersion)

	a.BootstrapInitialAdmin(
		app.Itrs.Bootstrap,
		cfg.InitialAdminEmail,
		cfg.InitialAdminPassword,
		*logger,
	)

	router := chi.NewRouter()
	webAuth := Chain(
		app.SessionManager.LoadAndSave,
		app.SessionManager.AuthCookiesMiddleWare,
		WebRequireLogin,
	)
	apiAuth := Chain(
		app.SessionManager.LoadAndSave,
		app.SessionManager.AuthBearerMiddleWare,
		app.SessionManager.AuthCookiesMiddleWare,
		ApiRequireLogin,
	)

	RouteWebPages(router, HomePageHandlerFunc(pageBuilder), webAuth)

	udb := userDashboard.New(pageBuilder, cfg.DefaultPageSize, app.Itrs.User.RenewToken,
		app.Itrs.User.ChangePassword, app.Itrs.Log.ListTasks, app.Itrs.Log.FindTask)
	udb.Route(router, webAuth)

	RouteAPI(router, *api.NewServer(&app.Itrs, *logger), apiAuth)
	RouteAPIDocs(router, APIDocsHandlerFunc(rt.APISpecsUrl, pageBuilder), webAuth)
	RouteAPISpecs(router)
	RouteStaticFiles(router)

	annotatorServer := an.NewServer(app.Annotator, pageBuilder, app.SessionManager)
	annotatorServer.Route(router, webAuth)

	collectionServer := clc.New(pageBuilder, cfg.DefaultPageSize,
		app.Itrs.Collection.Create, app.Itrs.Collection.List, app.Itrs.Collection.Update,
		app.Itrs.Collection.Delete, app.Itrs.Collection.Clone, app.Itrs.Collection.Find,
		app.Itrs.Group.List)
	collectionServer.Route(router, webAuth)

	imagesServer := im.New(
		pageBuilder,
		cfg.MaxArchiveMB,
		app.Itrs.Image.List,
		app.Itrs.Image.Delete,
		app.Itrs.Image.Find,
		app.Itrs.Image.IngestArchive,
	)
	imagesServer.Route(router, webAuth)

	adminPageBuilder := adm.NewPageBuilder(pageBuilder)
	adminUserServer := admusr.New(
		adminPageBuilder,
		app.Itrs.User,
		app.Itrs.Group,
		app.Itrs.Role,
		cfg.DefaultPageSize,
	)
	adminUserServer.Route(router, webAuth)
	adminGroupServer := admgrp.New(adminPageBuilder, app.Itrs.Group)
	adminGroupServer.Route(router, webAuth)
	adminRoleServer := admrl.New(adminPageBuilder, app.Itrs.Role)
	adminRoleServer.Route(router, webAuth)
	adminPolicyServer := admpl.New(adminPageBuilder, app.Itrs.Policy)
	adminPolicyServer.Route(router, webAuth)

	labelServer := lbl.New(pageBuilder, cfg.DefaultPageSize,
		app.Itrs.Label.Create, app.Itrs.Label.List, app.Itrs.Label.Update,
		app.Itrs.Label.Delete, app.Itrs.Label.Find)
	labelServer.Route(router, webAuth)

	notifier := wauth.MakeNotifierFromEnv(*logger)
	authServer := wauth.New(
		fmt.Sprintf("%v:%v", cfg.URL, port),
		basePageBuilder,
		*logger,
		app.SessionManager,
		notifier,
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
