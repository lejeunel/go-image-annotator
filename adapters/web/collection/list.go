package collection

import (
	_ "embed"
	"fmt"
	"io"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"

	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	e "github.com/lejeunel/go-image-annotator/adapters/web/error"
	pg "github.com/lejeunel/go-image-annotator/adapters/web/pagination"
	rt "github.com/lejeunel/go-image-annotator/routes"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/list"
)

//go:embed preamble.md
var preamble string

var listCollectionsFields = []string{
	"name",
	"description",
	"group",
	"profile",
	"created",
	"actions",
}

type MetaDeletePresenter struct {
	writer http.ResponseWriter
	b.RowURL
	okMessageFunc func(string) string
	htmx.ErrorPresenter
}

func NewMetaDeletePresenter(w http.ResponseWriter, u b.RowURL) MetaDeletePresenter {
	okMessageFunc := func(key string) string {
		return fmt.Sprintf("Successfully deleted key %v", key)
	}
	return MetaDeletePresenter{w, u, okMessageFunc, htmx.NewErrorPresenter("deleting meta-data", w)}
}

func (p MetaDeletePresenter) SuccessDeleteMetadata(key string) {
	htmx.NotifySuccessPayload(p.writer, "deleting meta-data", p.okMessageFunc(key))
}

func (p MetaDeletePresenter) SuccessReadMetadata(k string, v any) {
	b.RenderConfirmDeleteRow(len(listCollectionsFields),
		k, "meta-data", p.Url, p.writer)
}

func (s *Server) TableRow(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get(resourceUrlFieldName)
	s.RowURL.SetId(name)
	switch r.URL.Query().Get("mode") {
	case b.ModeEdit.String():
		p := NewEditPresenter(w, s.RowURL)
		s.FindItr.Execute(r.Context(), name, &p)
		s.ListProfileItr.Execute(r.Context(), &p)
		s.ListGroupItr.Execute(r.Context(), &p)
		p.Render()
	case b.ModeConfirmDelete.String():
		s.FindItr.Execute(r.Context(), name,
			NewDeletePresenter(w, s.RowURL))
	case b.ModeClone.String():
		s.FindItr.Execute(r.Context(), name,
			NewClonePresenter(w, s.RowURL))
	default:
		s.FindItr.Execute(r.Context(), name,
			NewViewPresenter(w, s.PageBuilder, s.RowURL))
	}
}

func (s *Server) List(w http.ResponseWriter, r *http.Request) {
	s.PageBuilder.SetUserIdentity(r.Context())
	s.ListCollectionItr.Execute(
		r.Context(),
		pa.PaginationParams{PageSize: s.DefaultPageSize, Page: pg.GetPageFromRequest(r)},
		NewListPresenter(w, s.PageBuilder, s.RowURL),
	)
}

type ListPresenter struct {
	b.PaginatedListBuilder
	b.RowURL
	Writer io.Writer
	e.ErrorPresenter
}

func NewListPresenter(w http.ResponseWriter, p b.PageBuilder, u b.RowURL) ListPresenter {
	p.SetTitle("Collections").
		SetHTMLTitle("Collections").
		SetActiveSection(cmp.CollectionsPageActive)
	p.AddMarkdownPreamble(preamble)
	b := b.NewPaginatedListBuilder(p, listCollectionsFields)
	return ListPresenter{b, u, w, e.NewErrorPresenter(w, p)}
}

func (p ListPresenter) SuccessListCollections(r list.Response) {
	p.SetPagination(r.Pagination, rt.CollectionsUrl)
	for _, c := range r.Collections {
		row := MakeRow(p.RowURL, c)
		p.AddRow(row)
	}
	p.AddCreationButton("Create", CreateCollectionFormUrl, createCollectionTargetDiv)
	p.Render(p.Writer)
}
