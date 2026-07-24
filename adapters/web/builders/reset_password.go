package builders

import (
	"io"

	st "github.com/lejeunel/go-image-annotator/adapters/web/styles"
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
	endpoint := rt.AddQueryParams(rt.ResetPasswordUrl, "token", b.token)
	return Div(Class("flex justify-center"),
		Span(
			Div(Class("flex justify-center text-gray-900 dark:text-white font-bold text-xl mt-4 mb-4"), Text(resetPasswordTitle)),
			Form(
				Attr("hx-post", endpoint.String()),
				Attr("hx-swap", "outerHTML"),
				Attr("hx-target", "this"),
				Class("bg-white dark:bg-gray-800 p-8 rounded-lg shadow-md w-80 mb-4"),
				Label(For("New password"), Text("New password"), Class(st.FormLabel)),
				Input(Type("password"), ID("password"), Name("password"), Required(),
					Class(st.FormInput)),
				Label(For("New password (repeat)"), Text("New password (repeat)"), Class(st.FormLabel)),
				Input(Type("password"), ID("password-repeat"), Name("password-repeat"), Required(),
					Class(st.FormInput)),
				Button(Type("submit"), Text("Submit"), Class(st.PasswordButtonClass)),
			),
		),
	)
}
func (b *ResetPasswordBuilder) SetToken(token string) {
	b.token = token
}
func (b *ResetPasswordBuilder) Render(w io.Writer) {
	b.BasePageBuilder.SetFrameContent(b.makeContent()).Render(w)
}

func NewResetPasswordBuilder(base BasePageBuilder) ResetPasswordBuilder {
	return ResetPasswordBuilder{BasePageBuilder: base}
}
