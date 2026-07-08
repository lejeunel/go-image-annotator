package builders

import (
	"fmt"
	"io"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

var OAuthButtonClass = "w-full whitespace-nowrap rounded-radius bg-surface-alt border border-surface-alt px-4 py-2 text-sm font-medium tracking-wide text-on-surface-strong transition hover:opacity-75 text-center focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-surface-alt active:opacity-100 active:outline-offset-0 disabled:opacity-75 disabled:cursor-not-allowed dark:bg-surface-dark-alt dark:border-surface-dark-alt dark:text-on-surface-dark-strong dark:focus-visible:outline-surface-dark-alt"
var PasswordButtonClass = "w-full rounded-radius bg-success border border-success px-4 py-2 text-center text-sm font-medium tracking-wide text-on-success hover:opacity-75 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary active:opacity-100 active:outline-offset-0 dark:bg-success-dark dark:border-success-dark dark:text-on-success-dark dark:focus-visible:outline-success-dark"
var signInTitle = "Sign in to ImageAnnotator"

type OAuthProvider struct {
	Name string
	URL  string
}

type LoginPageBuilder struct {
	BasePageBuilder
	OAuthProviders []OAuthProvider
}

func (b *LoginPageBuilder) AddOAuthProvider(name, url string) *LoginPageBuilder {
	b.OAuthProviders = append(b.OAuthProviders, OAuthProvider{Name: name, URL: url})
	return b
}
func (b *LoginPageBuilder) makeContent() Node {
	buttons := []Node{}
	for _, p := range b.OAuthProviders {
		button := A(Href(p.URL),
			Class(OAuthButtonClass),
			Text(fmt.Sprintf("Continue with %v", p.Name)))
		buttons = append(buttons, button)
	}
	return Div(Class("flex justify-center"),
		Span(
			Div(Class("flex justify-center text-gray-900 dark:text-white font-bold text-xl mt-4 mb-4"), Text(signInTitle)),
			Form(
				Action("/auth/login/password"),
				Method("POST"),
				Class("bg-white dark:bg-gray-800 p-8 rounded-lg shadow-md w-80 mb-4"),
				Label(For("email"), Text("Email"), Class("block text-sm font-medium text-gray-900 dark:text-white")),
				Input(Type("email"), ID("email"), Name("email"), Required(),
					Class("w-full px-3 py-2 mb-6 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent")),

				Div(Class("flex items-center justify-between"),
					Label(For("password"), Text("Password"), Class("block text-sm font-medium text-gray-900 dark:text-white")),
					A(Class("mt-2 block text-xs font-medium text-primary underline-offset-2 hover:underline focus:underline focus:outline-hidden dark:text-primary-dark"),
						Href("/forgot-password"), Text("Forgot password?")),
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
	)
}
func (b *LoginPageBuilder) Render(w io.Writer) {
	b.BasePageBuilder.SetContent(b.makeContent())
	b.Build().Render(w)

}

func NewLoginPageBuilder(base BasePageBuilder) *LoginPageBuilder {
	return &LoginPageBuilder{BasePageBuilder: base}
}
