package web

import (
	"io"
	"net/http"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"github.com/lejeunel/go-image-annotator/use-cases/label/list"
	. "maragu.dev/gomponents"
)

type ListLabelsPresenter struct {
	b.PageBuilder
	Writer io.Writer
	WebPageErrorPresenter
}

func (p ListLabelsPresenter) Success(r list.Response) {
	listBuilder := b.NewPaginatedListBuilder([]string{"name", "description", "actions"}, rt.Labels, r.Pagination)
	for _, l := range r.Labels {
		actions := b.NewActionsPanelBuilder()
		actions.SetEdit("/edit-url")
		actions.SetDelete("/delete-url")
		listBuilder.AddRow(Text(l.Name), Raw(l.Description), actions.Build())
	}
	p.PageBuilder.SetContent(listBuilder.Build(), nil)
	p.Render(p.Writer)
}

func (s *Server) ListLabels(w http.ResponseWriter, r *http.Request) {
	s.PageBuilder.SetUserIdentity(r.Context())
	s.Label.List.Execute(r.Context(),
		list.Request{PageSize: s.DefaultPageSize, Page: int64(GetPageFromRequest(r))},
		NewListLabelsPresenter(w, s.PageBuilder))
}

func NewListLabelsPresenter(w http.ResponseWriter, p b.PageBuilder) ListLabelsPresenter {
	return ListLabelsPresenter{
		*p.SetTitle("Labels").SetActive(b.LabelsPageActive),
		w, NewWebPageErrorPresenter(w)}
}
