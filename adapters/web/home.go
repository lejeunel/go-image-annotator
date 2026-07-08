package web

import (
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	. "maragu.dev/gomponents"
)

func MakeHomePage(pb b.PageBuilder) Node {
	pb.SetTitle("Home")
	pb.SetActive(b.HomePageActive)
	pb.SetContent(Text("Welcome."))
	return pb.Build()
}
