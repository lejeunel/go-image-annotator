package image

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	an "github.com/lejeunel/go-image-annotator/adapters/web/annotator"
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	tb "github.com/lejeunel/go-image-annotator/adapters/web/builders/table"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	ew "github.com/lejeunel/go-image-annotator/adapters/web/error"
	pg "github.com/lejeunel/go-image-annotator/adapters/web/pagination"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	rt "github.com/lejeunel/go-image-annotator/routes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
	find_im "github.com/lejeunel/go-image-annotator/use-cases/image/find"
	"github.com/lejeunel/go-image-annotator/use-cases/image/list"
	list_im "github.com/lejeunel/go-image-annotator/use-cases/image/list"
	. "maragu.dev/gomponents"
)

type ListImagesPresenter struct {
	b.PaginatedListBuilder
	io.Writer
	ew.ErrorPresenter
	collection string
}

var listImagesFields = []string{"id", "collection", "ingested", "n. annot.", "actions"}

func NewListImagesPresenter(w http.ResponseWriter, p b.PageBuilder, collection string) ListImagesPresenter {
	p.SetTitle(fmt.Sprintf("%v / Images", collection)).SetHTMLTitle("Images")
	b := b.NewPaginatedListBuilder(p, listImagesFields)
	return ListImagesPresenter{b, w, ew.NewErrorPresenter(w), collection}
}

func makeImageRow(image im.Image) tb.Row {
	link := rt.MakeAnnotateImageURL(an.AnnotateImage, image.Id.String(), image.Collection.Name)
	actions := b.NewActionsPanelBuilder()
	actions.SetConfirmDelete(rt.AddQueryParams(rt.Image, "id", image.Id.String(),
		"collection", image.Collection.Name,
		"mode", "confirm-delete"))
	row := tb.NewRow()
	row.AddCell(tb.NewCell(cmp.MakeTextLink(link, image.Id.String())))
	row.AddCell(tb.NewCell(Text(image.Collection.Name)))
	row.AddCell(tb.NewCell(Text(cmp.DateTimeToStr(image.Specs.IngestedAt))))
	row.AddCell(tb.NewCell(Text(strconv.Itoa(image.NumAnnotations()))))
	row.AddCell(tb.NewCell(actions.Build()))
	return row

}
func (p ListImagesPresenter) SuccessReadImage(image im.Image) {
	makeImageRow(image).Render(p.Writer)
}

func (p ListImagesPresenter) SuccessListImages(r list.Response) {

	baseURL := rt.AddQueryParams(rt.Images, "collection", p.collection)
	p.SetPagination(r.Pagination, baseURL.String())
	for _, im := range r.Images {
		p.AddRow(makeImageRow(im))
	}
	p.Render(p.Writer)
}

func (s *Server) List(w http.ResponseWriter, r *http.Request) {

	s.PageBuilder.SetUserIdentity(r.Context())
	collection := r.URL.Query().Get("collection")
	if collection == "" {
		s.PageBuilder.SetError(fmt.Errorf("parsing url to get collection name: %w", e.ErrURLParsing))
		s.PageBuilder.Render(w)
	}
	s.ListItr.Execute(list_im.Request{
		Filtering: im.Filtering{
			Collection: &collection},
		PaginationParams: pa.PaginationParams{
			PageSize: s.DefaultPageSize,
			Page:     pg.GetPageFromRequest(r)},
		Ordering: im.Ordering{IngestTime: true}},
		NewListImagesPresenter(w, s.PageBuilder, collection))
}

func (s *Server) TableRow(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	collection := r.URL.Query().Get("collection")
	switch r.URL.Query().Get("mode") {
	case b.ModeConfirmDelete.String():
		b.RenderConfirmDeleteRow(len(listImagesFields),
			id,
			"image",
			rt.AddQueryParams(rt.Image, "id", id, "collection", collection),
			w)
	default:
		s.FindItr.Execute(
			find_im.Request{
				ImageId:    id,
				Collection: collection,
			},
			NewListImagesPresenter(w, s.PageBuilder, collection))
	}
}
