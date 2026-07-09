package web

import (
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	"io"
	. "maragu.dev/gomponents"
)

func MakeHomePage(pb b.PageBuilder, w io.Writer) {
	pb.SetTitle("Home")
	pb.SetActive(b.HomePageActive)
	pb.SetContent(Text("Welcome."), nil)
	pb.Render(w)
}
