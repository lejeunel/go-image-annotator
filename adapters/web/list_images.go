package web

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	html "github.com/lejeunel/go-image-annotator/adapters/web/html"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/lejeunel/go-image-annotator/use-cases/image/list"
	list_im "github.com/lejeunel/go-image-annotator/use-cases/image/list"
	. "maragu.dev/gomponents"
)

type ListImagesPresenter struct {
	ListRenderer
}

func (p ListImagesPresenter) Success(r list.Response) {
	table := html.MyTable{Fields: []string{"id", "collection", "created", "n. annotations"}}
	for _, im := range r.Images {
		link := fmt.Sprintf("/image?id=%v&collection=%v", im.Id.String(), im.Collection.Name)
		table.Rows = append(table.Rows,
			html.MyTableRow{Values: []Node{html.MakeTextLink(link, im.Id.String()),
				Text(im.Collection.Name), Text(im.Specs.IngestedAt.Format(time.DateOnly)),
				Text(strconv.Itoa(im.NumAnnotations()))}})
	}
	p.RenderSuccess(table, r.Pagination, nil)
}

func (s *Server) ListImages(w http.ResponseWriter, r *http.Request) {

	s.PageBuilder.SetUserIdentityFromContext(r.Context())
	collection := r.URL.Query().Get("collection")
	if collection == "" {
		p := b.NewPageBuilder(s.APIPath).SetError(fmt.Errorf("parsing url to get collection name: %w", e.ErrURLParsing))
		p.Render(w)
	}
	s.Image.List.Execute(list_im.Request{PageSize: s.Image.DefaultPageSize,
		Page:           int64(GetPageFromRequest(r)),
		CollectionName: &collection},
		NewListImagesPresenter(w, *r.URL, s.PageBuilder))
}

func NewListImagesPresenter(w http.ResponseWriter, baseURL url.URL, pb b.PageBuilder) ListImagesPresenter {
	return ListImagesPresenter{
		ListRenderer: NewListRenderer(
			*pb.SetTitle("Images"),
			baseURL,
			b.NoPageActive, w),
	}
}
