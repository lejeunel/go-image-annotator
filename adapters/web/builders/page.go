package builders

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"maps"

	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	g "github.com/lejeunel/go-image-annotator/globals"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"github.com/yuin/goldmark"
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
	preamble            string
	content             Node
	BasePageBuilder
}

func NewPageBuilder(base BasePageBuilder, version g.Info) PageBuilder {
	return PageBuilder{BasePageBuilder: base, APIPath: rt.APIRootUrl, RepoURL: g.RepoURL, DocsURL: g.DocsURL,
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
	b.SidebarEntries = maps.Clone(b.SidebarEntries)
	for k, v := range b.SidebarEntries {
		v.IsActive = false
		if k == name {
			v.IsActive = true
		}
		b.SidebarEntries[k] = v

	}
	return b
}
func (b *PageBuilder) AddSidebarEntry(name, icon, url string, isActive bool) *PageBuilder {
	b.SidebarEntries = maps.Clone(b.SidebarEntries)
	b.SidebarEntries[name] = cmp.SidebarEntry{Label: name, Icon: icon, Url: url, IsActive: isActive}
	b.SidebarEntriesOrder = append(b.SidebarEntriesOrder, name)
	return b
}
func (b *PageBuilder) AddMarkdownPreamble(preamble string) *PageBuilder {

	md := goldmark.New()
	var buf bytes.Buffer
	if err := md.Convert([]byte(preamble), &buf); err != nil {
		panic(err)
	}

	b.preamble = buf.String()
	return b
}
func (b *PageBuilder) SetContent(content Node) *PageBuilder {
	b.content = content
	return b
}
func (b *PageBuilder) Render(w io.Writer) {
	if b.User == nil {
		b.BasePageBuilder.SetError(fmt.Errorf("current user has not been set"))
		b.BasePageBuilder.Render(w)
	}

	var header Node

	if b.Title != "" {
		header = Div(header, Div(Class("font-bold text-2xl"), Text(b.Title)))
	}
	if b.preamble != "" {
		header = Div(header, Article(Class("prose dark:prose-invert max-w-none"), Raw(b.preamble)))
	}

	if len(b.SidebarEntries) > 0 {
		var bufSidebar bytes.Buffer
		sidebar := cmp.NewSidebar(b.SidebarTitle)
		for _, n := range b.SidebarEntriesOrder {
			e, _ := b.SidebarEntries[n]
			sidebar.AddEntry(e.Label, e.Icon, e.Url, e.IsActive)
		}
		sidebar.Render(&bufSidebar)
		b.content = Div(Class("relative flex w-full flex-col"),
			Nav(Attr("x-cloak"),
				Class(`fixed left-0 top-14 z-20 flex h-svh w-60 shrink-0 flex-col border-r border-outline bg-surface-alt p-4 transition-transform duration-300
                      dark:border-outline-dark dark:bg-surface-dark-alt`), Raw(bufSidebar.String())),
			Div(Class("ml-60 px-4 py-18"), header, b.content),
		)
	} else {
		b.content = Div(Class("grow w-full px-4 py-18"), header, b.content)
	}

	b.BasePageBuilder.SetFrameContent(
		Group(
			[]Node{
				cmp.MakeNavBar(b.ActivePage, b.RepoURL, b.DocsURL, b.APIPath, *b.User, rt.UserDashboardUrl),
				b.content,
				cmp.MakeFooter(b.Version),
			},
		))
	b.BasePageBuilder.Render(w)
}
