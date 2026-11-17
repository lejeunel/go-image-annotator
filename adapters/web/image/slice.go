package image

import (
	"net/http"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	htmx "github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"github.com/lejeunel/go-image-annotator/use-cases/image/slice"
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
	cb := cmp.CopyURLButton()
	filteringCard := cmp.SpecCard{Fields: []cmp.SpecFields{
		{Name: "Filters", Value: filters},
		{Name: "Ordering", Value: ordering},
	}, Button: &cb}
	p.AddPreamble(filteringCard.Render())
	b := b.NewPaginatedListBuilder(p, listImagesFields)
	return SlicePresenter{b, w, htmx.NewErrorPresenter("querying", w), filters, ordering}
}

func (p SlicePresenter) SuccessSliceImages(r slice.Response) {
	baseURL := rt.MakeSliceURL(r.FilterStr, r.OrderStr)
	p.SetPagination(r.Pagination, baseURL.String())
	for _, im := range r.Images {
		p.AddRow(makeImageRow(im, MakeAnnotateImageURLFunc(im, p.FilterStr, p.OrderStr)))
	}
	p.Writer.Header().Set("HX-Push-Url", baseURL.String())

	p.Render(p.Writer)
}

func (s *Server) SliceFromForm(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}
	filters := r.FormValue(rt.FilterQueryArgName)
	ordering := r.FormValue(rt.OrderingQueryArgName)
	s.doSlice(w, r, filters, ordering)
}

func (s *Server) Slice(w http.ResponseWriter, r *http.Request) {
	filters := r.URL.Query().Get(rt.FilterQueryArgName)
	ordering := r.URL.Query().Get(rt.OrderingQueryArgName)
	s.doSlice(w, r, filters, ordering)
}

func (s *Server) doSlice(w http.ResponseWriter, r *http.Request, filters string, ordering string) {
	s.PageBuilder.SetUserIdentity(r.Context())

	s.SliceItr.Execute(slice.Request{
		FilterStr: filters,
		OrderStr:  ordering,
	}, NewSlicePresenter(w, s.PageBuilder, filters, ordering))
}
