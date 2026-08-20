package image

import (
	"fmt"
	"net/http"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	htmx "github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"github.com/lejeunel/go-image-annotator/use-cases/image/list"
)

type SlicePresenter struct {
	b.PaginatedListBuilder
	Writer http.ResponseWriter
	htmx.ErrorPresenter
	im.FilterStr
	im.OrderStr
}

func NewSlicePresenter(
	w http.ResponseWriter,
	p b.PageBuilder,
	filters im.FilterStr,
	ordering im.OrderStr,
) SlicePresenter {
	p.SetTitle("Slice").SetHTMLTitle("Slice")
	p.SetActiveSection(cmp.NoPageActive)
	p.AddMarkdownPreamble(fmt.Sprintf(`
**filters**: %v\
**ordering**: %v
`, filters, ordering))
	b := b.NewPaginatedListBuilder(p, listImagesFields)
	return SlicePresenter{b, w, htmx.NewErrorPresenter("querying", w), filters, ordering}
}

func (p SlicePresenter) SuccessListImages(r list.Response) {
	baseURL := rt.AddQueryParams(rt.SliceUrl,
		rt.FilterQueryArgName, r.FilterStr,
		rt.OrderingQueryArgName, r.OrderStr)
	p.SetPagination(r.Pagination, baseURL.String())
	for _, im := range r.Images {
		p.AddRow(makeImageRow(im, MakeAnnotateImageURLFunc(im, p.FilterStr, p.OrderStr)))
	}
	p.Writer.Header().Set("HX-Push-Url", baseURL.String())

	p.Render(p.Writer)
}

func (s *Server) Slice(w http.ResponseWriter, r *http.Request) {
	s.PageBuilder.SetUserIdentity(r.Context())
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}
	filters := r.FormValue(rt.FilterQueryArgName)
	ordering := r.FormValue(rt.OrderingQueryArgName)
	s.ListItr.Execute(list.Request{FilterStr: filters,
		OrderStr: ordering}, NewSlicePresenter(w, s.PageBuilder, filters, ordering))
}
