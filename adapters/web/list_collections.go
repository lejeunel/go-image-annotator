package web

import (
	"net/http"
	"net/url"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	html "github.com/lejeunel/go-image-annotator/adapters/web/html"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/list"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type ListCollectionsPresenter struct {
	ListRenderer
}

func (p ListCollectionsPresenter) SuccessListCollections(r list.Response) {
	table := html.MyTable{Fields: []string{"name", "description", "group", "created", "actions"}}
	for _, c := range r.Collections {
		var groupName string
		if c.Group == nil {
			groupName = "n/a"
		} else {
			groupName = c.Group.Name
		}

		actions := html.NewActionsPanel()
		actions.SetEdit("/edit-url")
		actions.SetDelete(rt.MakeDeleteCollectionURL(c.Name))
		table.AddRow(html.MakeTextLink(rt.MakeImagesURL(c.Name), c.Name),
			Raw(c.Description), Raw(groupName), Raw(cmp.DateTimeToStr(c.CreatedAt)), actions.Build())
	}
	button := cmp.MakeHTMXCreateButton("Create new collection", rt.CreateCollectionForm, createCollectionTargetDiv)
	preamble := Div(ID(createCollectionTargetDiv))
	p.RenderList(&preamble, table, r.Pagination, &button)
}

func (s *Server) ListCollections(w http.ResponseWriter, r *http.Request) {
	s.PageBuilder.SetUserIdentity(r.Context())
	s.Collection.List.Execute(r.Context(), list.Request{PageSize: s.DefaultPageSize, Page: int64(GetPageFromRequest(r))},
		NewListCollectionsPresenter(w, s.PageBuilder))
}
func NewListCollectionsPresenter(w http.ResponseWriter, p b.PageBuilder) ListCollectionsPresenter {
	baseURL, _ := url.Parse(rt.Collections)
	return ListCollectionsPresenter{
		ListRenderer: NewListRenderer(*p.SetTitle("Collections").SetActive(b.CollectionsPageActive), *baseURL,
			w),
	}
}
