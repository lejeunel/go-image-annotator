package builders

import (
	"io"

	s "github.com/lejeunel/go-image-annotator/adapters/web/styles"
	rt "github.com/lejeunel/go-image-annotator/routes"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

var forgotPasswordTitle = "Forgot Password"

type ForgotPasswordBuilder struct {
	BasePageBuilder
}

func (b *ForgotPasswordBuilder) Render(w io.Writer) {
	content := Div(Class("flex justify-center"),
		Span(
			Div(
				Class(
					"flex justify-center text-gray-900 dark:text-white font-bold text-xl mt-4 mb-4",
				),
				Text(forgotPasswordTitle),
			),
			Form(
				Attr("hx-post", rt.NotifyPasswordResetUrl),
				Attr("hx-swap", "outerHTML"),
				Attr("hx-target", "this"),
				Class("bg-white dark:bg-gray-800 p-8 rounded-lg shadow-md w-80 mb-4"),
				Label(
					For("email"),
					Text("Email"),
					Class("block text-sm font-medium text-gray-900 dark:text-white"),
				),
				Input(
					Type("email"),
					ID("email"),
					Name("email"),
					Required(),
					Class(
						"w-full px-3 py-2 mb-6 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent",
					),
				),
				Button(Type("submit"), Text("Submit"), Class(s.PasswordButtonClass)),
			),
		),
	)
	b.BasePageBuilder.SetFrameContent(content).Render(w)
}

func NewForgotPasswordBuilder(base BasePageBuilder) ForgotPasswordBuilder {
	base.SetHTMLTitle(forgotPasswordTitle)
	return ForgotPasswordBuilder{BasePageBuilder: base}
}
