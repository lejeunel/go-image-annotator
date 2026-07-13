package builders

import (
	"fmt"
	"io"

	s "github.com/lejeunel/go-image-annotator/adapters/web/styles"
	rt "github.com/lejeunel/go-image-annotator/routes"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

var signInTitle = "Sign in to ImageAnnotator"

type OAuthProvider struct {
	Name string
	URL  string
}

type LoginPageBuilder struct {
	BasePageBuilder
	OAuthProviders []OAuthProvider
}

func (b *LoginPageBuilder) AddOAuthProvider(provider, url string) *LoginPageBuilder {
	b.OAuthProviders = append(b.OAuthProviders, OAuthProvider{Name: provider, URL: url})
	return b
}
func (b *LoginPageBuilder) makeContent() Node {
	buttons := []Node{}
	for _, p := range b.OAuthProviders {
		button := A(Href(p.URL),
			Class(s.OAuthButtonClass),
			Text(fmt.Sprintf("Continue with %v", p.Name)))
		buttons = append(buttons, button)
	}
	return Div(Class("flex justify-center"),
		Span(
			Div(Class("flex justify-center text-gray-900 dark:text-white font-bold text-xl mt-4 mb-4"), Text(signInTitle)),
			Form(
				Action(rt.LoginWithPassword),
				Method("POST"),
				Class("bg-white dark:bg-gray-800 p-8 rounded-lg shadow-md w-80 mb-4"),
				Label(For("email"), Text("Email"), Class("block text-sm font-medium text-gray-900 dark:text-white")),
				Input(Type("email"), ID("email"), Name("email"), Required(),
					Class("w-full px-3 py-2 mb-6 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent")),

				Div(Class("flex items-center justify-between"),
					Label(For("password"), Text("Password"), Class("block text-sm font-medium text-gray-900 dark:text-white")),
					A(Class("mt-2 block text-xs font-medium text-primary underline-offset-2 hover:underline focus:underline focus:outline-hidden dark:text-primary-dark"),
						Href(rt.ForgotPasswordForm), Text("Forgot password?")),
				),
				Input(Type("password"), ID("password"), Name("password"), Required(),
					Class("w-full px-3 py-2 mb-6 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent")),
				Button(Type("submit"), Text("Login"), Class(s.PasswordButtonClass)),
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
	b.BasePageBuilder.SetContent(b.makeContent()).Render(w)
}

func NewLoginPageBuilder(base BasePageBuilder) LoginPageBuilder {
	return LoginPageBuilder{BasePageBuilder: base}
}
