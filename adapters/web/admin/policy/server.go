package policy

import (
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	pl "github.com/lejeunel/go-image-annotator/use-cases/policy"
)

type Server struct {
	Page b.PageBuilder
	Itrs pl.Interactors
}

func New(pb b.PageBuilder, itrs pl.Interactors) Server {
	pb.ActivateSidebarEntry(PageName)
	pb.SetHTMLTitle("Policies").SetTitle("Policies")
	return Server{pb, itrs}
}
