package admin

import (
	usr "github.com/lejeunel/go-image-annotator/adapters/web/admin/user"
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	"github.com/lejeunel/go-image-annotator/adapters/web/icons"
	rt "github.com/lejeunel/go-image-annotator/routes"
)

const (
	GroupPage  string = "Groups"
	RolePage   string = "Roles"
	PolicyPage string = "Policies"
)

func NewPageBuilder(pb b.PageBuilder) b.PageBuilder {

	pb.SetActiveSection(cmp.NoPageActive)
	pb.AddSidebarTitle("Admin")
	pb.AddSidebarEntry(usr.UserSidebarEntryName, icons.User, rt.AdminUsers, false)
	pb.AddSidebarEntry(GroupPage, icons.Group, rt.AdminGroups, false)
	pb.AddSidebarEntry(RolePage, icons.Rocket, rt.AdminRoles, false)
	pb.AddSidebarEntry(PolicyPage, icons.Shield, rt.AdminPolicies, false)
	return pb
}
