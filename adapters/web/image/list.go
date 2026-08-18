package image

import (
	"fmt"
	"io"
	"net/http"
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

func NewListImagesPresenter(
	w http.ResponseWriter,
	p b.PageBuilder,
	collection string,
) ListImagesPresenter {
	p.SetTitle(fmt.Sprintf("%v / Images", collection)).SetHTMLTitle("Images")
	p.SetActiveSection(cmp.NoPageActive)
	b := b.NewPaginatedListBuilder(p, listImagesFields)
	return ListImagesPresenter{b, w, ew.NewErrorPresenter(w), collection}
}

func (p ListImagesPresenter) SuccessReadImage(image im.Image) {
	makeImageRow(image, BaseAnnotateImageURLFunc).Render(p.Writer)
}

func (p ListImagesPresenter) SuccessListImages(r list.Response) {
	baseURL := rt.AddQueryParams(rt.ImagesUrl, rt.CollectionArgName, p.collection)
	p.SetPagination(r.Pagination, baseURL.String())
	ingestUrl := rt.AddQueryParams(ingestPanelUrl, rt.CollectionArgName, p.collection)
	p.PaginatedListBuilder.AddCreationButton("Ingest", ingestUrl.String(), ingestTargetDiv)
	for _, im := range r.Images {
		p.AddRow(makeImageRow(im, BaseAnnotateImageURLFunc))
	}
	p.Render(p.Writer)
}

func MakeAnnotateImageURL(baseURL, imageId, collection string) string {
	return fmt.Sprintf("%v?id=%v&collection=%v", baseURL, imageId, collection)
}

func makeImageRow(image im.Image, urlFunc ImageURLFunc) tb.Row {
	link := urlFunc(image)
	actions := b.NewActionsPanelBuilder()
	actions.SetConfirmDelete(rt.AddQueryParams(ImageRow, "id", image.Id.String(),
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

func (s *Server) List(w http.ResponseWriter, r *http.Request) {
	s.PageBuilder.SetUserIdentity(r.Context())
	collection := r.URL.Query().Get("collection")
	if collection == "" {
		s.PageBuilder.SetError(
			fmt.Errorf("parsing url to get collection name: %w", e.ErrURLParsing),
		)
		s.PageBuilder.Render(w)
	}
	s.ListItr.Execute(list_im.Request{
		FilterStr:        fmt.Sprintf("collection=\"%v\"", collection),
		PaginationParams: pa.PaginationParams{Page: pg.GetPageFromRequest(r)},
		OrderStr:         "ingested_at:asc",
	},
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
			rt.AddQueryParams(ImageRow, "id", id, "collection", collection),
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
