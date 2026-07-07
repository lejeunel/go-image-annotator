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

var OAuthButtonClass = "w-full whitespace-nowrap rounded-radius bg-surface-alt border border-surface-alt px-4 py-2 text-sm font-medium tracking-wide text-on-surface-strong transition hover:opacity-75 text-center focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-surface-alt active:opacity-100 active:outline-offset-0 disabled:opacity-75 disabled:cursor-not-allowed dark:bg-surface-dark-alt dark:border-surface-dark-alt dark:text-on-surface-dark-strong dark:focus-visible:outline-surface-dark-alt"
var PasswordButtonClass = "w-full rounded-radius bg-success border border-success px-4 py-2 text-center text-sm font-medium tracking-wide text-on-success hover:opacity-75 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary active:opacity-100 active:outline-offset-0 dark:bg-success-dark dark:border-success-dark dark:text-on-success-dark dark:focus-visible:outline-success-dark"
var title = "Sign in to ImageAnnotator"

func (b *LoginPageBuilder) Build() Node {
	buttons := []Node{}
	for _, p := range b.OAuthProviders {
		button := A(Href(p.URL),
			Class(OAuthButtonClass),
			Text(fmt.Sprintf("Continue with %v", p.Name)))
		buttons = append(buttons, button)
	}
	return Doctype(HTML(
		Head(
			Title(title),
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
				Div(Class("flex justify-center"),
					Span(
						Div(Class("font-bold text-xl mb-4"), Text(title)),
						Form(
							Action("/auth/login/password"),
							Method("POST"),
							Class("bg-white p-8 rounded-lg shadow-md w-80 mb-4"),
							Label(For("email"), Text("Email"), Class("block text-sm font-medium text-gray-700")),
							Input(Type("email"), ID("email"), Name("email"), Required(),
								Class("w-full px-3 py-2 mb-6 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent")),

							Div(Class("flex items-center justify-between"),
								Label(For("password"), Text("Password"), Class("block text-sm font-medium text-gray-700")),
								A(Class("mt-2 block text-xs font-medium text-primary underline-offset-2 hover:underline focus:underline focus:outline-hidden dark:text-primary-dark"),
									Href("/auth/reset-password"), Text("Forgot password?")),
							),
							Input(Type("password"), ID("password"), Name("password"), Required(),
								Class("w-full px-3 py-2 mb-6 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent")),
							Button(Type("submit"), Text("Login"), Class(PasswordButtonClass)),
						),
						Div(Class("flex items-center gap-3 my-6 mb-4"),
							Div(Class("h-px flex-1 bg-gray-300")),
							Span(Class("text-xs"), Text("or")),
							Div(Class("h-px flex-1 bg-gray-300"))),
						Div(Class("flex justify-center"),
							Group(buttons),
						),
					),
				),
			),
		),
	))

}

func (b *LoginPageBuilder) Render(w io.Writer) {
	b.Build().Render(w)

}

func NewLoginPageBuilder() *LoginPageBuilder {
	return &LoginPageBuilder{}
}
