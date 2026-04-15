package web

import (
	"fmt"
	"net/http"
	"net/url"

	e "github.com/lejeunel/go-image-annotator-v2/shared/errors"
	html "github.com/lejeunel/go-image-annotator-v2/shared/html"
	n "github.com/lejeunel/go-image-annotator-v2/shared/navigation"
	"github.com/lejeunel/go-image-annotator-v2/use-cases/image/list"
	list_im "github.com/lejeunel/go-image-annotator-v2/use-cases/image/list"
	. "maragu.dev/gomponents"
)

type ListImagesPresenter struct {
	ListRenderer
}

func (p ListImagesPresenter) Success(r list.Response) {
	table := html.PaginationTable{Fields: []string{"id", "collection"}}
	for _, im := range r.Images {
		link := fmt.Sprintf("/image?id=%v&collection=%v", im.Id.String(), im.Collection.Name)
		table.Rows = append(table.Rows,
			html.PaginationTableRow{Values: []Node{html.MakeTextLink(link, im.Id.String()),
				Text(im.Collection.Name)}})
	}
	p.RenderSuccess(table, r.Pagination)
}

func (s *Server) ListImages(w http.ResponseWriter, r *http.Request) {

	collection := r.URL.Query().Get("collection")
	if collection == "" {
		p := html.NewPageBuilder().SetError(fmt.Errorf("parsing url to get collection name: %w", e.ErrURLParsing))
		p.Render(w)
	}
	s.Image.List.Execute(list_im.Request{PageSize: s.Image.DefaultPageSize,
		Page:           int64(GetPageFromRequest(r)),
		CollectionName: &collection},
		NewListImagesPresenter(w, *r.URL))
}

func NewListImagesPresenter(w http.ResponseWriter, baseURL url.URL) ListImagesPresenter {
	return ListImagesPresenter{
		ListRenderer: NewListRenderer("Images", baseURL,
			n.NoPageActive, w),
	}
}
