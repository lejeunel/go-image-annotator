package web

import (
	html "github.com/lejeunel/go-image-annotator/shared/html"
	n "github.com/lejeunel/go-image-annotator/shared/navigation"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
)

func MakeHomePage() Node {

	p := html.NewTitledPageBuilder("Home")
	p.SetContent(Div(Text("This is a sentence with a "), html.MakeTextLink("#", "link")))
	p.SetActive(n.HomePageActive)
	return p.Build()
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	MakeHomePage().Render(w)

}
