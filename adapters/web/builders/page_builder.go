package builders

import (
	"bytes"
	"context"
	"fmt"
	"maps"

	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	g "github.com/lejeunel/go-image-annotator/globals"
	rt "github.com/lejeunel/go-image-annotator/routes"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type PageBuilder struct {
	APIPath             string
	RepoURL             string
	DocsURL             string
	Version             g.Info
	ActivePage          cmp.ActivePage
	User                *u.User
	SidebarTitle        string
	SidebarEntries      map[string]cmp.SidebarEntry
	SidebarEntriesOrder []string
	BasePageBuilder
}

func NewPageBuilder(base BasePageBuilder, version g.Info) PageBuilder {
	return PageBuilder{BasePageBuilder: base, APIPath: rt.APIRoot, RepoURL: g.RepoURL, DocsURL: g.DocsURL,
		Version: version, SidebarEntries: make(map[string]cmp.SidebarEntry)}
}
func (b *PageBuilder) SetHTMLTitle(title string) *PageBuilder {
	b.BasePageBuilder.SetHTMLTitle(title)
	return b
}
func (b *PageBuilder) SetTitle(title string) *PageBuilder {
	b.BasePageBuilder.SetTitle(title)
	return b
}
func (b *PageBuilder) SetActiveSection(a cmp.ActivePage) *PageBuilder {
	b.ActivePage = a
	return b
}
func (b *PageBuilder) SetUserIdentity(ctx context.Context) *PageBuilder {
	id := u.IdentityFromContext(ctx)
	b.User = id
	return b
}
func (b *PageBuilder) AddSidebarTitle(title string) *PageBuilder {
	b.SidebarTitle = title
	return b
}
func (b *PageBuilder) ActivateSidebarEntry(name string) *PageBuilder {
	m, ok := b.SidebarEntries[name]
	if !ok {
		return b
	}
	m.IsActive = true
	b.SidebarEntries[name] = m
	return b
}

func (b *PageBuilder) AddSidebarEntry(name, icon, url string, isActive bool) *PageBuilder {
	b.SidebarEntries = maps.Clone(b.SidebarEntries)
	b.SidebarEntries[name] = cmp.SidebarEntry{Label: name, Icon: icon, Url: url, IsActive: isActive}
	b.SidebarEntriesOrder = append(b.SidebarEntriesOrder, name)
	return b
}
func (b *PageBuilder) SetContent(content Node) *PageBuilder {
	if b.User == nil {
		b.BasePageBuilder.SetError(fmt.Errorf("current user has not been set"))
		return b
	}

	if b.Title != "" {
		content = Div(Div(Class("font-bold text-xl"), Text(b.Title)),
			content)

	}

	if len(b.SidebarEntries) > 0 {
		var bufSidebar bytes.Buffer
		sidebar := cmp.NewSidebar(b.SidebarTitle)
		for _, n := range b.SidebarEntriesOrder {
			e, _ := b.SidebarEntries[n]
			sidebar.AddEntry(e.Label, e.Icon, e.Url, e.IsActive)
		}
		sidebar.Render(&bufSidebar)
		content = Div(Class("relative flex w-full flex-col"),
			Nav(Attr("x-cloak"),
				Class(`fixed left-0 top-14 z-20 flex h-svh w-60 shrink-0 flex-col border-r border-outline bg-surface-alt p-4 transition-transform duration-300
                      dark:border-outline-dark dark:bg-surface-dark-alt`), Raw(bufSidebar.String())),
			Div(Class("ml-60"), content),
		)
	}

	b.BasePageBuilder.SetFrameContent(
		Group(
			[]Node{
				cmp.MakeNavBar(b.ActivePage, b.RepoURL, b.DocsURL, b.APIPath, *b.User),
				Div(Class("grow w-full px-1 px-4 py-18"),
					content,
				),
				cmp.MakeFooter(b.Version),
			},
		))
	return b
}
