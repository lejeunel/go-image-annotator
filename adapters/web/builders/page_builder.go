package builders

import (
	"context"
	"fmt"

	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	g "github.com/lejeunel/go-image-annotator/globals"
	rt "github.com/lejeunel/go-image-annotator/routes"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type PageBuilder struct {
	APIPath    string
	RepoURL    string
	DocsURL    string
	Version    g.Info
	ActivePage cmp.ActivePage
	User       *u.User
	BasePageBuilder
}

func NewPageBuilder(base BasePageBuilder, version g.Info) *PageBuilder {
	return &PageBuilder{BasePageBuilder: base, APIPath: rt.APIRoot, RepoURL: g.RepoURL, DocsURL: g.DocsURL, Version: version}
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
func (b *PageBuilder) SetContent(content Node) *PageBuilder {
	if b.User == nil {
		b.BasePageBuilder.SetError(fmt.Errorf("current user has not been set"))
		return b
	}

	b.BasePageBuilder.SetContent(
		Group(
			[]Node{
				cmp.MakeNavBar(b.ActivePage, b.RepoURL, b.DocsURL, b.APIPath, *b.User),
				Div(Class("grow w-full px-1 md:px-2 lg:px-4 py-10 md:py-20"),
					Div(Class("font-bold text-xl"), Text(b.Title)),
					content,
				),
				cmp.MakeFooter(b.Version),
			},
		))
	return b
}
