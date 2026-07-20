package builders

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"io"

	"github.com/lejeunel/go-image-annotator/adapters/web/icons"
	rt "github.com/lejeunel/go-image-annotator/routes"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

//go:embed templates/admin.html
var AdminHTML string

type ActivePage int

const (
	UserPageActive ActivePage = iota
	GroupPageActive
	RolePageActive
	PolicyPageActive
)

type SidebarItem struct {
	URL    string
	Label  string
	Icon   template.HTML
	Active bool
}

type AdminPageData struct {
	Content template.HTML
	Sidebar template.HTML
}

type AdminPageBuilder struct {
	PageBuilder
	active ActivePage
}

func makeSidebarEntry(icon, label, url string, isActive bool) Node {
	aClass := "flex items-center rounded-radius gap-2 px-2 py-1.5 text-sm font-medium text-on-surface underline-offset-2 hover:bg-primary/5 hover:text-on-surface-strong focus-visible:underline focus:outline-hidden dark:text-on-surface-dark dark:hover:bg-primary-dark/5 dark:hover:text-on-surface-dark-strong"
	if isActive {
		aClass = "flex items-center rounded-radius gap-2 bg-primary/10 px-2 py-1.5 text-sm font-medium text-on-surface-strong underline-offset-2 focus-visible:underline focus:outline-hidden dark:bg-primary-dark/10 dark:text-on-surface-dark-strong"
	}
	return A(Href(url), Class(aClass), Raw(icon), Span(Text(label)))
}

func (b *AdminPageBuilder) SetActive(a ActivePage) *AdminPageBuilder {
	b.active = a
	return b
}

func (b *AdminPageBuilder) renderContent(w io.Writer) {
	Text("Page").Render(w)
}
func (b *AdminPageBuilder) renderSidebar(w io.Writer) {
	items := []Node{
		makeSidebarEntry(icons.User, "Users", rt.AdminUsers, b.active == UserPageActive),
		makeSidebarEntry(icons.Group, "Groups", rt.AdminGroups, b.active == GroupPageActive),
		makeSidebarEntry(icons.Rocket, "Roles", rt.AdminRoles, b.active == RolePageActive),
		makeSidebarEntry(icons.Shield, "Policies", rt.AdminPolicies, b.active == PolicyPageActive),
	}
	nodes := Div(Class("flex flex-col gap-2 overflow-y-auto pb-6"), Group(items))
	nodes.Render(w)
}

func (b *AdminPageBuilder) Build() *AdminPageBuilder {
	var content bytes.Buffer
	b.renderContent(&content)
	var sidebar bytes.Buffer
	b.renderSidebar(&sidebar)
	baseTemplate, _ := template.New("").Parse(AdminHTML)

	var out bytes.Buffer
	if err := baseTemplate.ExecuteTemplate(&out, "admin",
		AdminPageData{template.HTML(content.String()),
			template.HTML(sidebar.String())}); err != nil {
		b.SetError(fmt.Errorf("executing admin template: %w", err))
		return b
	}

	b.SetContent(Raw(out.String()))
	return b
}

func NewAdminPageBuilder(b PageBuilder) AdminPageBuilder {
	return AdminPageBuilder{PageBuilder: b}
}
