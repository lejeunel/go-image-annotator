package builders

import (
	"fmt"
	"io"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type OAuthProvider struct {
	Name string
	URL  string
}

type LoginPageBuilder struct {
	OAuthProviders []OAuthProvider
}

func (b *LoginPageBuilder) AddOAuthProvider(name, url string) *LoginPageBuilder {
	b.OAuthProviders = append(b.OAuthProviders, OAuthProvider{Name: name, URL: url})
	return b
}

func (b *LoginPageBuilder) Build() Node {
	buttons := []Node{}
	for _, p := range b.OAuthProviders {
		button := A(Href(p.URL),
			Class("rounded-radius bg-primary border border-primary px-4 py-2 text-center text-sm font-medium tracking-wide text-on-primary hover:opacity-75 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary active:opacity-100 active:outline-offset-0 dark:bg-primary-dark dark:border-primary-dark dark:text-on-primary-dark dark:focus-visible:outline-primary-dark"),
			Text(fmt.Sprintf("Login with %v", p.Name)))
		buttons = append(buttons, button)
	}
	return Doctype(HTML(
		Head(
			Title("Login"),
			Meta(Charset("utf-8")),
			Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
			Link(
				Rel("stylesheet"),
				Href("/static/styles.css"),
			),
			Link(Rel("stylesheet"), Href("https://fonts.googleapis.com/css2?family=Roboto&display=swap")),
		),
		Body(
			Class("bg-white text-gray-900 dark:bg-gray-900 dark:text-white"),
			Div(Class("grow w-full px-1 md:px-2 lg:px-4 py-10 md:py-20"),
				Div(Class("font-bold text-xl mb-4"), Text("Login")),
				Group(buttons)),
		),
	))

}

func (b *LoginPageBuilder) Render(w io.Writer) {
	b.Build().Render(w)

}

func NewLoginPageBuilder() *LoginPageBuilder {
	return &LoginPageBuilder{}
}
