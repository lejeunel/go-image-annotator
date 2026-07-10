package web

import (
	"io"
	"net/http"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	html "github.com/lejeunel/go-image-annotator/adapters/web/html"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/list"
	. "maragu.dev/gomponents"
)

type ListCollectionsPresenter struct {
	b.PageBuilder
	Writer io.Writer
	WebPageErrorPresenter
}

func (p ListCollectionsPresenter) SuccessListCollections(r list.Response) {
	listBuilder := b.NewPaginatedListBuilder(
		[]string{"name", "description", "group", "created", "actions"},
		rt.Collections,
		r.Pagination)
	for _, c := range r.Collections {
		var groupName string
		if c.Group == nil {
			groupName = "n/a"
		} else {
			groupName = c.Group.Name
		}

		actions := b.NewActionsPanelBuilder()
		actions.SetEdit("/edit-url")
		actions.SetDelete(rt.MakeDeleteCollectionURL(c.Name))
		listBuilder.AddRow(
			html.MakeTextLink(rt.MakeImagesURL(c.Name), c.Name),
			Raw(c.Description), Raw(groupName), Raw(cmp.DateTimeToStr(c.CreatedAt)), actions.Build())
	}
	listBuilder.AddCreationButton("Create new collection", rt.CreateCollectionForm, createCollectionTargetDiv)
	p.PageBuilder.SetContent(listBuilder.Build(), nil)
	p.Render(p.Writer)
}

func (s *Server) ListCollections(w http.ResponseWriter, r *http.Request) {
	s.PageBuilder.SetUserIdentity(r.Context())
	s.Collection.List.Execute(r.Context(), list.Request{PageSize: s.DefaultPageSize, Page: int64(GetPageFromRequest(r))},
		NewListCollectionsPresenter(w, s.PageBuilder))
}
func NewListCollectionsPresenter(w http.ResponseWriter, p b.PageBuilder) ListCollectionsPresenter {
	return ListCollectionsPresenter{*p.SetTitle("Collections").SetActive(b.CollectionsPageActive), w,
		NewWebPageErrorPresenter(w)}
}
