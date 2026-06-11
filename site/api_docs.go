package site

import (
	"context"
	html "github.com/lejeunel/go-image-annotator/shared/html"
	n "github.com/lejeunel/go-image-annotator/shared/navigation"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func APIDocsPage(ctx context.Context, specsPath string, apiPath string) Node {
	p := html.NewPageBuilder(apiPath)
	p.AddScripts(html.APIDocsLib())
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
	p.SetActive(n.APIDocsPageActive)
	return p.Build()
}
