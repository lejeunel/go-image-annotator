package web

import (
	html "github.com/lejeunel/go-image-annotator/shared/html"
	n "github.com/lejeunel/go-image-annotator/shared/navigation"
	. "maragu.dev/gomponents"
	"net/http"
)

func MakeHomePage(b html.PageBuilder) Node {
	b.SetTitle("Home")
	b.SetActive(n.HomePageActive)
	return b.Build()
}

func HomePageHandlerFunc(b html.PageBuilder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b.SetUserIdentityFromContext(r.Context())
		MakeHomePage(b).Render(w)
	}
}
