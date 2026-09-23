package dashboard

import (
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	"github.com/lejeunel/go-image-annotator/adapters/web/icons"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	rt "github.com/lejeunel/go-image-annotator/routes"
	ft "github.com/lejeunel/go-image-annotator/use-cases/log/find"
	lt "github.com/lejeunel/go-image-annotator/use-cases/log/list"
	cpw "github.com/lejeunel/go-image-annotator/use-cases/user/change-password"
	rat "github.com/lejeunel/go-image-annotator/use-cases/user/renew-access-token"
)

type Server struct {
	b.PageBuilder
	RenewAPITokenItr  rat.Interactor
	ChangePasswordItr cpw.Interactor
	ListTasksItr      lt.Interactor
	FindTaskItr       ft.Interactor
	DefaultPageSize   int
	policies          auth.Policies
}

func New(
	pb b.PageBuilder,
	defaultPageSize int,
	i rat.Interactor,
	c cpw.Interactor,
	lt lt.Interactor,
	ft ft.Interactor,
	policies auth.Policies,
) Server {
	pb.AddSidebarEntry(ProfilePageName, icons.Info, rt.DashboardUrl, false)
	pb.AddSidebarEntry(CredentialsPageName, icons.Key, CredentialsUrl, false)
	pb.AddSidebarEntry(LogsPageName, icons.Notepad, rt.ListTasksUrl, false)
	return Server{pb, i, c, lt, ft, defaultPageSize, policies}
}
