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
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	rt "github.com/lejeunel/go-image-annotator/routes"
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
	p.SetTitle("Collections").SetHTMLTitle("Collections").SetActiveSection(cmp.CollectionsPageActive)
	b := b.NewPaginatedListBuilder(p, listCollectionsFields)
	return ListCollectionsPresenter{b, w, e.NewErrorPresenter(w)}
}
func (p ListCollectionsPresenter) SuccessListCollections(r list.Response) {
	p.SetPagination(r.Pagination, rt.Collections)
	for _, c := range r.Collections {
		row := MakeListCollectionRow(c)
		p.AddRow(row)
	}
	p.AddCreationButton("Create", CreateCollectionForm, createCollectionTargetDiv)
	p.PaginatedListBuilder.AddMarkdownPreamble(preamble)
	p.Render(p.Writer)
}

type CollectionPresenter struct {
	Writer io.Writer
	e.WebPageErrorPresenter
	successFindCollection func(clc.Collection)
}

func NewCollectionPresenter(w http.ResponseWriter, mode string) CollectionPresenter {
	p := CollectionPresenter{Writer: w, WebPageErrorPresenter: e.NewErrorPresenter(w)}
	switch mode {
	case "edit":
		p.successFindCollection = p.renderEditForm
	case "confirm-delete":
		p.successFindCollection = p.renderConfirmDelete
	default:
		p.successFindCollection = p.renderView
	}
	return p
}
func (p CollectionPresenter) SuccessFindCollection(c clc.Collection) {
	p.successFindCollection(c)
}
func (p *CollectionPresenter) renderEditForm(c clc.Collection) {
	endpoint := rt.AddQueryParams(Collection, "name", c.Name)
	b := bf.NewHTMXInlineFormBuilder(c.Name, len(listCollectionsFields), endpoint)
	b.AddTitle(fmt.Sprintf("Editing %v", c.Name))
	b.AddTextField("name", "Name", "name", bf.WithRequired(), bf.WithDefault(c.Name))
	b.AddTextField("description", "Description", "description", bf.WithDefault(c.Description))
	b.Render(p.Writer)
}
func (p *CollectionPresenter) renderConfirmDelete(c clc.Collection) {
	b.RenderConfirmDeleteRow(len(listCollectionsFields),
		c.Name,
		"collection",
		rt.AddQueryParams(Collection, "name", c.Name),
		rt.AddQueryParams(Collection, "name", c.Name, "mode", "view"),
		p.Writer)
}
func (p *CollectionPresenter) renderView(c clc.Collection) {
	MakeListCollectionRow(c).Render(p.Writer)
}

func MakeListCollectionRow(c clc.Collection) tb.Row {
	var groupName string
	if c.Group == nil {
		groupName = "n/a"
	} else {
		groupName = c.Group.Name
	}

	actions := b.NewActionsPanelBuilder()
	actions.SetEdit(rt.AddQueryParams(Collection,
		"name", c.Name, "description", c.Description, "mode", "edit"))
	actions.SetConfirmDelete(rt.AddQueryParams(Collection, "name", c.Name,
		"mode", "confirm-delete"))

	row := tb.NewRow()
	row.AddCell(tb.NewCell(cmp.MakeTextLink(rt.MakeImagesURL(c.Name), c.Name)))
	row.AddCell(tb.NewCell(Text(c.Description)))
	row.AddCell(tb.NewCell(Text(groupName)))
	row.AddCell(tb.NewCell(Text(cmp.DateTimeToStr(c.CreatedAt))))
	row.AddCell(tb.NewCell(actions.Build()))
	return row

}
