package web

import (
	"context"
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
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

func RegisterAPIDocs(mux *http.ServeMux, specsPath, docsPath string, p b.PageBuilder) {
	mux.HandleFunc(docsPath,
		func(w http.ResponseWriter, r *http.Request) {
			if err := APIDocsPage(r.Context(), specsPath, *p.SetActive(b.APIDocsPageActive)).Render(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
}
