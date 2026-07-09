package web

import (
	"context"
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	"io"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func APIDocsLib() Node {
	return Script(Src("https://unpkg.com/@stoplight/elements/web-components.min.js"))

}

func APIDocsPage(ctx context.Context, specsPath string, p b.PageBuilder, w io.Writer) {
	p.AddScripts(APIDocsLib())
	p.SetUserIdentityFromContext(ctx)
	p.SetTitle("API Docs")
	p.SetContent(Div(Class("spotlight "),
		El("elements-api",
			Attr("apiDescriptionUrl", specsPath),
			// This includes session cookies when calling endpoints with the
			// "try-it" button
			Attr("tryItCredentialsPolicy", "include"),
			Attr("router", "hash"),
			Attr("layout", "sidebar"),
		)), nil)
	p.SetActive(b.APIDocsPageActive)
	p.Render(w)
}
