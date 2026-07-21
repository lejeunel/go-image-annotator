package admin

import (
	_ "embed"

	"github.com/go-chi/chi/v5"
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	"github.com/lejeunel/go-image-annotator/adapters/web/icons"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"net/http"
)

const (
	UserPage   string = "Users"
	GroupPage  string = "Groups"
	RolePage   string = "Roles"
	PolicyPage string = "Policies"
)

type Server struct {
	b.SideBarPageBuilder
}

func (s *Server) Users(w http.ResponseWriter, r *http.Request) {
	s.SetUserIdentity(r.Context())
	s.SetActiveSidebarItem(UserPage)
	s.Build().Render(w)
}

func (s *Server) Route(r chi.Router, mws ...func(http.Handler) http.Handler) {

	r.Group(func(r chi.Router) {
		r.Use(mws...)
		r.Get(rt.Admin, s.Users)
		r.Get(rt.AdminUsers, s.Users)
	})
}

func New(pb b.PageBuilder) Server {
	sbp := b.SideBarPageBuilder{PageBuilder: pb, Sidebar: cmp.NewSidebar("Admin")}
	sbp.Sidebar.AddEntry(UserPage, icons.User, rt.AdminUsers)
	sbp.Sidebar.AddEntry(GroupPage, icons.Group, rt.AdminGroups)
	sbp.Sidebar.AddEntry(RolePage, icons.Rocket, rt.AdminRoles)
	sbp.Sidebar.AddEntry(PolicyPage, icons.Shield, rt.AdminPolicies)
	sbp.SetActiveSection(cmp.NoPageActive)
	return Server{sbp}
}
