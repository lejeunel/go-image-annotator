package web

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	html "github.com/lejeunel/go-image-annotator/adapters/web/html"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	rt "github.com/lejeunel/go-image-annotator/routes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/lejeunel/go-image-annotator/use-cases/image/list"
	list_im "github.com/lejeunel/go-image-annotator/use-cases/image/list"
	. "maragu.dev/gomponents"
)

type ListImagesPresenter struct {
	ListRenderer
}

func (p ListImagesPresenter) SuccessListImages(r list.Response) {
	table := html.MyTable{Fields: []string{"id", "collection", "ingested", "n. annot.", "actions"}}
	for _, im := range r.Images {
		link := rt.MakeImageURL(im.Id.String(), im.Collection.Name)
		actions := html.NewActionsPanel()
		actions.SetDelete("/delete-url")
		table.AddRow(
			html.MakeTextLink(link, im.Id.String()),
			Text(im.Collection.Name), Text(cmp.DateTimeToStr(im.Specs.IngestedAt)),
			Text(strconv.Itoa(im.NumAnnotations())), actions.Build())
	}
	p.RenderList(nil, table, r.Pagination, nil)
}

func (s *Server) ListImages(w http.ResponseWriter, r *http.Request) {

	s.PageBuilder.SetUserIdentity(r.Context())
	collection := r.URL.Query().Get("collection")
	if collection == "" {
		s.PageBuilder.SetError(fmt.Errorf("parsing url to get collection name: %w", e.ErrURLParsing))
		s.PageBuilder.Render(w)
	}
	s.Image.List.Execute(list_im.Request{
		FilteringParams: im.FilteringParams{
			Collection: &collection,
			PageSize:   s.DefaultPageSize,
			Page:       int64(GetPageFromRequest(r))},
		OrderingParams: im.OrderingParams{IngestTime: true}},
		NewListImagesPresenter(w, *r.URL, s.PageBuilder))
}

func NewListImagesPresenter(w http.ResponseWriter, baseURL url.URL, pb b.PageBuilder) ListImagesPresenter {
	return ListImagesPresenter{
		ListRenderer: NewListRenderer(
			*pb.SetTitle("Images"),
			baseURL,
			w),
	}
}
