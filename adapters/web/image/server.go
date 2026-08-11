package image

import (
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	"github.com/lejeunel/go-image-annotator/use-cases/image/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/image/find"
	ia "github.com/lejeunel/go-image-annotator/use-cases/image/ingest-archive"
	"github.com/lejeunel/go-image-annotator/use-cases/image/list"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type Server struct {
	b.PageBuilder
	DefaultPageSize  int
	ListItr          list.Interactor
	DeleteItr        delete.Interactor
	FindItr          find.Interactor
	IngestArchiveItr ia.Interactor
}

func CodeHighlightingLibs() []Node {
	var scripts []Node
	scripts = append(scripts, Link(Href("/static/prism.css"), Rel("stylesheet")))
	scripts = append(scripts, Script(Src("/static/prism.js")))
	scripts = append(scripts, Script(Src("/static/prism-python.js")))
	return scripts
}

func New(
	pb b.PageBuilder, defaultPageSize int,
	l list.Interactor, d delete.Interactor, f find.Interactor,
	i ia.Interactor,
) Server {
	pb.AddScripts(CodeHighlightingLibs()...)
	return Server{pb,
		defaultPageSize, l, d, f, i}
}
