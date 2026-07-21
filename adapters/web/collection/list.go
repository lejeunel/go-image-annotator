package collection

import (
	"fmt"
	"io"
	"net/http"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	bf "github.com/lejeunel/go-image-annotator/adapters/web/builders/form"
	tb "github.com/lejeunel/go-image-annotator/adapters/web/builders/table"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	e "github.com/lejeunel/go-image-annotator/adapters/web/error"
	pg "github.com/lejeunel/go-image-annotator/adapters/web/pagination"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	rt "github.com/lejeunel/go-image-annotator/routes"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/list"
	. "maragu.dev/gomponents"
)

var listCollectionsFields = []string{"name", "description", "group", "created", "actions"}

type ListCollectionsPresenter struct {
	b.PaginatedListBuilder
	Writer io.Writer
	e.WebPageErrorPresenter
}

func NewListCollectionsPresenter(w http.ResponseWriter, p b.PageBuilder) ListCollectionsPresenter {
	p.SetTitle("Collections").SetActiveSection(cmp.CollectionsPageActive)
	b := b.NewPaginatedListBuilder(p, listCollectionsFields)
	return ListCollectionsPresenter{b, w, e.NewErrorPresenter(w)}
}
func (p ListCollectionsPresenter) SuccessFindCollection(c clc.Collection) {
	MakeListCollectionRow(c).Render(p.Writer)
}
func (p ListCollectionsPresenter) SuccessListCollections(r list.Response) {
	p.SetPagination(r.Pagination, rt.Collections)
	for _, c := range r.Collections {
		row := MakeListCollectionRow(c)
		p.AddRow(row)
	}
	p.AddCreationButton("Create new collection", rt.CreateCollectionForm, createCollectionTargetDiv)
	p.Build().Render(p.Writer)
}

func (s *Server) TableRow(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	switch r.URL.Query().Get("mode") {
	case "edit":
		endpoint := rt.AddQueryParams(rt.Collection, "name", name)
		b := bf.NewHTMXInlineFormBuilder(name, len(listCollectionsFields), endpoint)
		b.AddTitle(fmt.Sprintf("Editing %v", name))
		b.AddTextField("name", "Name", "name", bf.WithRequired(), bf.WithDefault(name))
		b.AddTextField("description", "Description", "description", bf.WithDefault(r.URL.Query().Get("description")))
		b.Render(w)
	case "confirm-delete":
		b.RenderConfirmDeleteRow(len(listCollectionsFields),
			name,
			"collection",
			rt.AddQueryParams(rt.Collection, "name", name),
			rt.AddQueryParams(rt.Collection, "name", name, "mode", "view"),
			w)
	default:
		s.FindItr.Execute(r.Context(),
			r.URL.Query().Get("name"),
			NewListCollectionsPresenter(w, s.PageBuilder))
	}

}
func (s *Server) CreateForm(w http.ResponseWriter, r *http.Request) {
	b := bf.NewHTMXCreateFormBuilder(rt.Collection, createCollectionTargetDiv)
	b.AddTitle("Create a new collection")
	b.AddTextField("name", "Name", "name", bf.WithRequired())
	b.AddTextField("description", "Description", "description")
	b.Render(w)
}
func (s *Server) List(w http.ResponseWriter, r *http.Request) {
	s.PageBuilder.SetUserIdentity(r.Context())
	s.ListItr.Execute(r.Context(), pa.PaginationParams{PageSize: s.DefaultPageSize, Page: pg.GetPageFromRequest(r)},
		NewListCollectionsPresenter(w, s.PageBuilder))
}

func MakeListCollectionRow(c clc.Collection) tb.Row {
	var groupName string
	if c.Group == nil {
		groupName = "n/a"
	} else {
		groupName = c.Group.Name
	}

	actions := b.NewActionsPanelBuilder()
	actions.SetEdit(rt.AddQueryParams(rt.Collection,
		"name", c.Name, "description", c.Description, "mode", "edit"))
	actions.SetConfirmDelete(rt.AddQueryParams(rt.Collection, "name", c.Name,
		"mode", "confirm-delete"))

	row := tb.NewRow()
	row.AddCell(tb.NewCell(cmp.MakeTextLink(rt.MakeImagesURL(c.Name), c.Name)))
	row.AddCell(tb.NewCell(Text(c.Description)))
	row.AddCell(tb.NewCell(Text(groupName)))
	row.AddCell(tb.NewCell(Text(cmp.DateTimeToStr(c.CreatedAt))))
	row.AddCell(tb.NewCell(actions.Build()))
	return row

}
