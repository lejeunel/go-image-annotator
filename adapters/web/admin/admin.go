package admin

import (
	grp "github.com/lejeunel/go-image-annotator/adapters/web/admin/group"
	pl "github.com/lejeunel/go-image-annotator/adapters/web/admin/policy"
	rl "github.com/lejeunel/go-image-annotator/adapters/web/admin/role"
	usr "github.com/lejeunel/go-image-annotator/adapters/web/admin/user"
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	"github.com/lejeunel/go-image-annotator/adapters/web/icons"
	rt "github.com/lejeunel/go-image-annotator/routes"
)

const (
	RolePage   string = "Roles"
	PolicyPage string = "Policies"
)

func NewPageBuilder(pb b.PageBuilder) b.PageBuilder {
	pb.SetActiveSection(cmp.NoPageActive)
	pb.AddSidebarTitle("Admin")
	pb.AddSidebarEntry(usr.PageName, icons.User, rt.AdminUsersUrl, false)
	pb.AddSidebarEntry(grp.PageName, icons.Group, rt.AdminGroupsUrl, false)
	pb.AddSidebarEntry(rl.PageName, icons.Rocket, rt.AdminRolesUrl, false)
	pb.AddSidebarEntry(pl.PageName, icons.Shield, rt.AdminPoliciesUrl, false)
	return pb
}
