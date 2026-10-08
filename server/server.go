package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"syscall"

	api "github.com/lejeunel/go-image-annotator/adapters/api/server"
	userDashboard "github.com/lejeunel/go-image-annotator/adapters/web/dashboard"
	rt "github.com/lejeunel/go-image-annotator/routes"

	adm "github.com/lejeunel/go-image-annotator/adapters/web/admin"
	admgrp "github.com/lejeunel/go-image-annotator/adapters/web/admin/group"
	admrl "github.com/lejeunel/go-image-annotator/adapters/web/admin/role"
	admusr "github.com/lejeunel/go-image-annotator/adapters/web/admin/user"
	an "github.com/lejeunel/go-image-annotator/adapters/web/annotator"
	wauth "github.com/lejeunel/go-image-annotator/adapters/web/auth"
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	clc "github.com/lejeunel/go-image-annotator/adapters/web/collection"
	home "github.com/lejeunel/go-image-annotator/adapters/web/home"
	im "github.com/lejeunel/go-image-annotator/adapters/web/image"
	lbl "github.com/lejeunel/go-image-annotator/adapters/web/label"
	pr "github.com/lejeunel/go-image-annotator/adapters/web/profile"
	a "github.com/lejeunel/go-image-annotator/app"
	"github.com/lejeunel/go-image-annotator/app/sqlite"
	g "github.com/lejeunel/go-image-annotator/globals"

	"github.com/go-chi/chi/v5"
)

// Make initializes the root handler and listens on the given port.
func Make(port int) (http.Handler, *slog.Logger) {
	app := sqlite.NewAppFromEnv()
	app.DBManager.Init(context.Background())
	currentVersion := g.Info{Version: g.Version, Date: g.Date}
	basePageBuilder := b.NewBasePageBuilder()

	queryDocs := b.QueryDocs{
		Filtering:         app.ImageFilterDocumenter.DescribeFields(),
		FilteringExamples: app.ImageFilterDocumenter.Examples(),
		Ordering:          app.ImageSortDocumenter.DescribeOrderingFields(),
		OrderingExamples:  app.ImageSortDocumenter.Examples(),
	}
	pageBuilder := b.NewPageBuilder(basePageBuilder, currentVersion, queryDocs)

	a.BootstrapInitialAdmin(
		app.Itrs.Bootstrap,
		app.Config.InitialAdminEmail,
		app.Config.InitialAdminPassword,
		app.Logger,
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

	RouteWebPages(router, home.HandlerFunc(pageBuilder), webAuth)

	udb := userDashboard.New(pageBuilder, app.Config.DefaultPageSize, app.Itrs.User.RenewToken,
		app.Itrs.User.ChangePassword, app.Itrs.Log.ListTasks, app.Itrs.Log.FindTask)

	udb.Route(router, webAuth)

	RouteAPI(router, *api.NewServer(&app.Itrs, app.Logger), apiAuth)
	RouteAPIDocs(router, ApiDocsHandlerFunc(rt.APISpecsUrl, pageBuilder), webAuth)
	RouteAPISpecs(router)
	RouteStaticFiles(router)

	annotatorServer := an.NewServer(app.Annotator, pageBuilder, app.SessionManager)
	annotatorServer.Route(router, webAuth)

	collectionServer := clc.New(pageBuilder, app.Config.DefaultPageSize,
		app.Itrs.Collection.Create, app.Itrs.Collection.List, app.Itrs.Collection.Update,
		app.Itrs.Collection.Delete, app.Itrs.Collection.Clone, app.Itrs.Collection.Find,
		app.Itrs.Group.List, app.Itrs.Profile.ListAll)
	collectionServer.Route(router, webAuth)

	imagesServer := im.New(
		pageBuilder,
		app.Config.MaxArchiveMB,
		app.Itrs.Image.Slice,
		app.Itrs.Image.List,
		app.Itrs.Image.Delete,
		app.Itrs.Image.Find,
		app.Itrs.Image.IngestArchive,
	)
	imagesServer.Route(router, webAuth)

	profilesServer := pr.New(
		pageBuilder,
		app.Config.DefaultPageSize,
		app.Itrs.Profile.Create,
		app.Itrs.Profile.List,
		app.Itrs.Annotation.PickLabel,
		app.Itrs.Profile.Update,
		app.Itrs.Profile.Delete,
		app.Itrs.Profile.Find,
	)
	profilesServer.Route(router, webAuth)

	adminPageBuilder := adm.NewPageBuilder(pageBuilder)
	adminUserServer := admusr.New(
		adminPageBuilder,
		app.Itrs.User,
		app.Itrs.Group,
		app.Itrs.Role,
		app.Config.DefaultPageSize,
	)
	adminUserServer.Route(router, webAuth)
	adminGroupServer := admgrp.New(adminPageBuilder, app.Itrs.Group)
	adminGroupServer.Route(router, webAuth)
	adminRoleServer := admrl.New(adminPageBuilder, app.Itrs.Role, app)
	adminRoleServer.Route(router, webAuth)

	labelServer := lbl.New(pageBuilder, app.Config.DefaultPageSize,
		app.Itrs.Label.Create, app.Itrs.Label.List, app.Itrs.Label.Update,
		app.Itrs.Label.Delete, app.Itrs.Label.Find)
	labelServer.Route(router, webAuth)

	notifier := wauth.MakeNotifierFromEnv(app.Logger)
	authServer := wauth.New(
		fmt.Sprintf("%v:%v", app.Config.URL, port),
		app.Config,
		basePageBuilder,
		app.Logger,
		app.SessionManager,
		notifier,
		app.Itrs.User.RequestForgottenPassword,
		app.Itrs.User.ResetForgottenPassword)
	authServer.Route(router,
		app.SessionManager.LoadAndSave)

	return router, &app.Logger
}

func Serve(port int) {
	handler, logger := Make(port)
	addr := fmt.Sprintf(":%d", port)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		switch {
		case errors.Is(err, syscall.EADDRINUSE):
			logger.Error("port is already in use", "port", port)
		case errors.Is(err, syscall.EACCES):
			logger.Error("permission denied binding port", "port", port)
		default:
			logger.Error("cannot listen", "port", port, "err", err)
		}
		os.Exit(1)
	}

	logger.Info("started server", "port", port)

	if err := http.Serve(ln, handler); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server error", "err", err)
		os.Exit(1)
	}
}
