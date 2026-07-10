package builders

import (
	"context"
	"fmt"

	u "github.com/lejeunel/go-image-annotator/entities/user"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type PageBuilder struct {
	APIPath    string
	RepoURL    string
	DocsURL    string
	ActivePage ActivePage
	User       *u.User
	BasePageBuilder
}

func (b *PageBuilder) SetPane(pane Node) *PageBuilder {
	b.BasePageBuilder.SetPane(pane)
	return b
}

func (b *PageBuilder) SetTitle(title string) *PageBuilder {
	b.BasePageBuilder.SetTitle(title)
	return b
}
func (b *PageBuilder) SetActive(a ActivePage) *PageBuilder {
	b.ActivePage = a
	return b
}
func (b *PageBuilder) SetUserIdentity(ctx context.Context) *PageBuilder {
	id := u.IdentityFromContext(ctx)
	b.User = id
	return b
}
func (b *PageBuilder) SetContent(c Node, pane *Node) *PageBuilder {
	if b.User == nil {
		b.BasePageBuilder.SetError(fmt.Errorf("current user has not been set"))
		return b
	}

	var title Node
	if pane != nil {
		title = Span(Class("w-full inline-flex items-center justify-between"), Text(b.Title), *pane)
	} else {
		title = Text(b.Title)
	}
	b.BasePageBuilder.SetContent(
		Group(
			[]Node{MakeNavBar(b.ActivePage, b.RepoURL, b.DocsURL, b.APIPath, *b.User),
				Div(Class("grow w-full px-1 md:px-2 lg:px-4 py-10 md:py-20"),
					Div(Class("font-bold text-xl"), title),
					c)},
		))
	return b
}

func NewPageBuilder(base BasePageBuilder, apiPrefix, repoURL, docsURL string) *PageBuilder {
	return &PageBuilder{BasePageBuilder: base, APIPath: apiPrefix, RepoURL: repoURL, DocsURL: docsURL}
}
