package web

import (
	"io"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"

	. "maragu.dev/gomponents"
)

func MakeHomePage(pb b.PageBuilder, w io.Writer) {
	pb.SetTitle("Home")
	pb.SetActiveSection(cmp.HomePageActive)
	pb.SetContent(Text("Welcome."))
	pb.Render(w)
}
