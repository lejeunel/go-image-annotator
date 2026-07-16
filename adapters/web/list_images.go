package web

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	tb "github.com/lejeunel/go-image-annotator/adapters/web/builders/table"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	ew "github.com/lejeunel/go-image-annotator/adapters/web/error"
	pg "github.com/lejeunel/go-image-annotator/adapters/web/pagination"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	rt "github.com/lejeunel/go-image-annotator/routes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
	"github.com/lejeunel/go-image-annotator/use-cases/image/list"
	list_im "github.com/lejeunel/go-image-annotator/use-cases/image/list"
	. "maragu.dev/gomponents"
)

type ListImagesPresenter struct {
	b.PaginatedListBuilder
	io.Writer
	ew.WebPageErrorPresenter
	collection string
}

var listImagesFields = []string{"id", "collection", "ingested", "n. annot.", "actions"}

func NewListImagesPresenter(w http.ResponseWriter, p b.PageBuilder, collection string) ListImagesPresenter {
	p.SetTitle("Image")
	b := b.NewPaginatedListBuilder(p, listImagesFields)
	return ListImagesPresenter{b, w, ew.NewErrorPresenter(w), collection}
}

func (p ListImagesPresenter) SuccessListImages(r list.Response) {

	baseURL := rt.AddQueryParams(rt.Images, "collection", p.collection)
	p.SetPagination(r.Pagination, baseURL.String())
	for _, im := range r.Images {
		link := rt.MakeImageURL(im.Id.String(), im.Collection.Name)
		actions := b.NewActionsPanelBuilder()

		deleteURL, _ := url.Parse("/edit-url")
		actions.SetConfirmDelete(*deleteURL)
		row := tb.NewRow()
		row.AddCell(tb.NewCell(cmp.MakeTextLink(link, im.Id.String())))
		row.AddCell(tb.NewCell(Text(im.Collection.Name)))
		row.AddCell(tb.NewCell(Text(cmp.DateTimeToStr(im.Specs.IngestedAt))))
		row.AddCell(tb.NewCell(Text(strconv.Itoa(im.NumAnnotations()))))
		row.AddCell(tb.NewCell(actions.Build()))
		p.AddRow(row)
	}
	p.Build().Render(p.Writer)
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
			Collection: &collection},
		PaginationParams: pa.PaginationParams{
			PageSize: s.DefaultPageSize,
			Page:     pg.GetPageFromRequest(r)},
		OrderingParams: im.OrderingParams{IngestTime: true}},
		NewListImagesPresenter(w, s.PageBuilder, collection))
}
