package generic

import (
	"context"
	au "datahub/app/authorizer"
	gp "maragu.dev/gomponents"
	gh "maragu.dev/gomponents/html"
)

type GenericViewer struct {
	SignOutURL       string
	IdentityProvider au.IdentityProvider
}

func NewGenericViewer(signOutURL string, identityProvider au.IdentityProvider) *GenericViewer {
	return &GenericViewer{SignOutURL: signOutURL,
		IdentityProvider: identityProvider}
}

func (v *GenericViewer) BasePage(ctx context.Context, title string, head, body gp.Node) gp.Node {
	return HTMLPage("DataHub / "+title,
		head,
		gh.Div(gh.Class("min-h-screen flex flex-col justify-between"),
			NavBar(ctx, v.IdentityProvider, v.SignOutURL),
			gh.Div(gh.Class("grow w-full px-2 md:px-4 lg:px-8 py-2 md:py-4"),
				gh.H1(gh.Class("text-4xl font-bold text-gray-900 py-2"), gp.Text(title)),
				body,
			),
		))
}

func (v *GenericViewer) OopsPage(ctx context.Context, message string) gp.Node {
	return v.BasePage(ctx, "Oops...", Head(), gh.P(gp.Text(message)))
}
