package builders

import (
	"context"
	"io"

	u "github.com/lejeunel/go-image-annotator/entities/user"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type PageBuilder struct {
	Title      string
	APIPath    string
	RepoURL    string
	DocsURL    string
	scripts    []Node
	ActivePage ActivePage
	User       *u.User
	Content    Node
}

func (b *PageBuilder) AddScripts(scripts ...Node) *PageBuilder {
	for _, s := range scripts {
		b.scripts = append(b.scripts, s)
	}
	return b
}
func (b *PageBuilder) SetTitle(title string) *PageBuilder {
	b.Title = title
	return b
}

func (b *PageBuilder) SetActive(a ActivePage) *PageBuilder {
	b.ActivePage = a
	return b
}
func (b *PageBuilder) SetUserIdentityFromContext(ctx context.Context) *PageBuilder {
	id := u.IdentityFromContext(ctx)
	b.User = id
	return b
}
func (b *PageBuilder) SetContent(c Node) *PageBuilder {
	b.Content = c
	return b
}

func (b *PageBuilder) SetError(err error) *PageBuilder {
	b.Title = "Oops!"
	b.Content = Text(err.Error())
	return b
}

func (b *PageBuilder) Build() Node {
	b.scripts = append(b.scripts, BaseLibs()...)
	return Doctype(HTML(
		Attr("x-data", `{
					darkMode: false,

					init() {
						this.darkMode = localStorage.getItem('dark') === 'true'
						document.documentElement.classList.toggle('dark', this.darkMode)
					},
					toggleDark() {
						this.darkMode = !this.darkMode;
						localStorage.setItem('dark', this.darkMode);
						document.documentElement.classList.toggle('dark', this.darkMode)
					}}`),
		Attr("x-init", "init()"),
		Attr("x-bind:class", "{ 'dark': darkMode }"),
		Head(
			Title(b.Title),
			Meta(Charset("utf-8")),
			Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
			Script(Raw(`
				if (localStorage.getItem('dark') === 'true') {
					document.documentElement.classList.add('dark');
				}
			`)),
			Link(
				Rel("stylesheet"),
				Href("/static/styles.css"),
			),
			Link(Rel("stylesheet"), Href("https://fonts.googleapis.com/css2?family=Roboto&display=swap")),
		),
		Body(
			Class("bg-white text-gray-900 dark:bg-gray-900 dark:text-white"),
			MakeNavBar(b.ActivePage, b.RepoURL, b.DocsURL, b.APIPath, b.User),
			Div(Class("grow w-full px-1 md:px-2 lg:px-4 py-10 md:py-20"),
				Div(Class("font-bold text-xl"), Text(b.Title)),
				b.Content),
			Group(b.scripts),
		),
	))
}

func (b *PageBuilder) Render(w io.Writer) {
	b.Build().Render(w)

}

func NewPageBuilder(apiPrefix, repoURL, docsURL string) *PageBuilder {
	return &PageBuilder{APIPath: apiPrefix, RepoURL: repoURL, DocsURL: docsURL}
}
