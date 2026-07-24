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
	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/list"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/update"
	. "maragu.dev/gomponents"
)

var listCollectionsFields = []string{"name", "description", "group", "created", "actions"}

type ListPresenter struct {
	b.PaginatedListBuilder
	Writer io.Writer
	e.ErrorPresenter
}

func NewListPresenter(w http.ResponseWriter, p b.PageBuilder) ListPresenter {
	p.SetTitle("Collections").SetHTMLTitle("Collections").SetActiveSection(cmp.CollectionsPageActive)
	b := b.NewPaginatedListBuilder(p, listCollectionsFields)
	return ListPresenter{b, w, e.NewErrorPresenter(w)}
}
func (p ListPresenter) SuccessListCollections(r list.Response) {
	p.SetPagination(r.Pagination, rt.Collections)
	for _, c := range r.Collections {
		row := MakeRow(c)
		p.AddRow(row)
	}
	p.AddCreationButton("Create", CreateCollectionForm, createCollectionTargetDiv)
	p.PaginatedListBuilder.AddMarkdownPreamble(preamble)
	p.Render(p.Writer)
}

type RowPresenter struct {
	Writer io.Writer
	e.ErrorPresenter
	successFindCollection func(clc.Collection)
}

func NewCollectionPresenter(w http.ResponseWriter, mode string) RowPresenter {
	p := RowPresenter{Writer: w, ErrorPresenter: e.NewErrorPresenter(w)}
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
func (p RowPresenter) SuccessFindCollection(c clc.Collection) {
	p.successFindCollection(c)
}
func (p *RowPresenter) renderEditForm(c clc.Collection) {
	endpoint := rt.AddQueryParams(Collection, "name", c.Name)
	b := bf.NewHTMXInlineFormBuilder(c.Name, len(listCollectionsFields), endpoint)
	b.AddTitle(fmt.Sprintf("Editing %v", c.Name))
	b.AddTextField("name", "Name", "name", bf.WithRequired(), bf.WithDefault(c.Name))
	b.AddTextField("description", "Description", "description", bf.WithDefault(c.Description))
	b.Render(p.Writer)
}
func (p *RowPresenter) renderConfirmDelete(c clc.Collection) {
	b.RenderConfirmDeleteRow(len(listCollectionsFields),
		c.Name,
		"collection",
		rt.AddQueryParams(Collection, "name", c.Name),
		rt.AddQueryParams(Collection, "name", c.Name, "mode", "view"),
		p.Writer)
}
func (p *RowPresenter) renderView(c clc.Collection) {
	MakeRow(c).Render(p.Writer)
}

type EditCollectionPresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(update.Response) string
	htmx.ErrorPresenter
}

func NewEditCollectionPresenter(w http.ResponseWriter) EditCollectionPresenter {
	task := "Updating collection"
	okMessageFunc := func(r update.Response) string {
		return "Successfully updated collection"
	}
	return EditCollectionPresenter{w, task, okMessageFunc, htmx.NewErrorPresenter(task, w)}
}

func (p EditCollectionPresenter) SuccessUpdateCollection(r update.Response) {
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, p.okMessageFunc(r))
}

func MakeRow(c clc.Collection) tb.Row {
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
