package builders

import (
	"io"

	s "github.com/lejeunel/go-image-annotator/adapters/web/styles"
	rt "github.com/lejeunel/go-image-annotator/routes"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

var resetPasswordTitle = "Resetting password"

type ResetPasswordBuilder struct {
	BasePageBuilder
	token string
}

func (b *ResetPasswordBuilder) makeContent() Node {
	endpoint := rt.AddQueryParams(rt.ResetPassword, "token", b.token)
	return Div(Class("flex justify-center"),
		Span(
			Div(Class("flex justify-center text-gray-900 dark:text-white font-bold text-xl mt-4 mb-4"), Text(resetPasswordTitle)),
			Form(
				Attr("hx-post", endpoint.String()),
				Attr("hx-swap", "outerHTML"),
				Attr("hx-target", "this"),
				Class("bg-white dark:bg-gray-800 p-8 rounded-lg shadow-md w-80 mb-4"),
				Label(For("New password"), Text("New password"), Class("block text-sm font-medium text-gray-900 dark:text-white")),
				Input(Type("password"), ID("password"), Name("password"), Required(),
					Class("w-full px-3 py-2 mb-6 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent")),
				Label(For("New password (repeat)"), Text("New password (repeat)"), Class("block text-sm font-medium text-gray-900 dark:text-white")),
				Input(Type("password"), ID("password-repeat"), Name("password-repeat"), Required(),
					Class("w-full px-3 py-2 mb-6 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent")),
				Button(Type("submit"), Text("Submit"), Class(s.PasswordButtonClass)),
			),
		),
	)
}
func (b *ResetPasswordBuilder) SetToken(token string) {
	b.token = token
}
func (b *ResetPasswordBuilder) Render(w io.Writer) {
	b.BasePageBuilder.SetContent(b.makeContent()).Render(w)
}

func NewResetPasswordBuilder(base BasePageBuilder) ResetPasswordBuilder {
	return ResetPasswordBuilder{BasePageBuilder: base}
}
