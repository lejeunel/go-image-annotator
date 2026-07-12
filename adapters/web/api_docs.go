package web

import (
	"context"
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	"io"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func APIDocsLib() Node {
	return Script(Src("/static/stoplight.js"))

}

func APIDocsPage(ctx context.Context, specsPath string, p b.PageBuilder, w io.Writer) {
	p.AddScripts(APIDocsLib())
	p.SetUserIdentity(ctx)
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
	p.SetActive(cmp.APIDocsPageActive)
	p.Render(w)
}
