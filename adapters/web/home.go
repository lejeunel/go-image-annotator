package web

import (
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	. "maragu.dev/gomponents"
	"net/http"
)

func MakeHomePage(pb b.PageBuilder) Node {
	pb.SetTitle("Home")
	pb.SetActive(b.HomePageActive)
	pb.SetContent(Text("Welcome."))
	return pb.Build()
}

func HomePageHandlerFunc(pb b.PageBuilder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pb.SetUserIdentityFromContext(r.Context())
		MakeHomePage(pb).Render(w)
	}
}
