package web

import (
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	"io"
	. "maragu.dev/gomponents"
)

func MakeHomePage(pb b.PageBuilder, w io.Writer) {
	pb.SetTitle("Home")
	pb.SetActive(cmp.HomePageActive)
	pb.SetContent(Text("Welcome."))
	pb.Render(w)
}
