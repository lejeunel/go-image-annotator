package web

import (
	"context"
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func APIDocsLib() Node {
	return Script(Src("https://unpkg.com/@stoplight/elements/web-components.min.js"))

}

func APIDocsPage(ctx context.Context, specsPath string, p b.PageBuilder) Node {
	p.AddScripts(APIDocsLib())
	p.SetUserIdentityFromContext(ctx)
	p.SetContent(Div(Class("spotlight "),
		El("elements-api",
			Attr("apiDescriptionUrl", specsPath),
			// This includes session cookies when calling endpoints with the
			// "try-it" button
			Attr("tryItCredentialsPolicy", "include"),
			Attr("router", "hash"),
			Attr("layout", "sidebar"),
		)))
	p.SetActive(b.APIDocsPageActive)
	return p.Build()
}
