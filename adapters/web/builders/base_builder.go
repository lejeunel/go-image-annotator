package builders

import (
	"fmt"
	"io"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type BasePageBuilder struct {
	Title   string
	scripts []Node
	Error   error
	Content Node
}

func (b *BasePageBuilder) AddScripts(scripts ...Node) *BasePageBuilder {
	for _, s := range scripts {
		b.scripts = append(b.scripts, s)
	}
	return b
}
func (b *BasePageBuilder) SetTitle(title string) *BasePageBuilder {
	b.Title = title
	return b
}
func (b *BasePageBuilder) SetContent(c Node) *BasePageBuilder {
	b.Content = c
	return b
}
func (b *BasePageBuilder) SetError(err error) *BasePageBuilder {
	b.Error = err
	return b
}
func (b *BasePageBuilder) Build() Node {
	if b.Error != nil {
		return Text(b.Error.Error())
	}
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
			Group(b.scripts),
			Raw(fmt.Sprintf("<title>%v</title>", b.Title)),
		),
		Body(
			Class("bg-white text-gray-900 dark:bg-gray-900 dark:text-white"),
			b.Content),
	))
}
func (b *BasePageBuilder) Render(w io.Writer) {
	b.Build().Render(w)

}

func NewBasePageBuilder() BasePageBuilder {
	return BasePageBuilder{}
}
