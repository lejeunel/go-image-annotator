package policy

import (
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
)

type Server struct {
	Page b.PageBuilder
}

func New(pb b.PageBuilder) Server {
	pb.ActivateSidebarEntry(PageName)
	pb.SetHTMLTitle("Policies").SetTitle("Policies")
	return Server{pb}
}
