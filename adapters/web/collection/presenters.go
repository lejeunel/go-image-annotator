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
	b.RowURL
	Writer io.Writer
	e.ErrorPresenter
}

func NewListPresenter(w http.ResponseWriter, p b.PageBuilder, u b.RowURL) ListPresenter {
	p.SetTitle("Collections").SetHTMLTitle("Collections").SetActiveSection(cmp.CollectionsPageActive)
	b := b.NewPaginatedListBuilder(p, listCollectionsFields)
	return ListPresenter{b, u, w, e.NewErrorPresenter(w)}
}
func (p ListPresenter) SuccessListCollections(r list.Response) {
	p.SetPagination(r.Pagination, rt.CollectionsUrl)
	for _, c := range r.Collections {
		row := MakeRow(p.RowURL, c)
		p.AddRow(row)
	}
	p.AddCreationButton("Create", CreateCollectionFormUrl, createCollectionTargetDiv)
	p.PaginatedListBuilder.AddMarkdownPreamble(preamble)
	p.Render(p.Writer)
}

type ViewPresenter struct {
	Writer io.Writer
	b.RowURL
	e.ErrorPresenter
}

func NewViewPresenter(w http.ResponseWriter, u b.RowURL) ViewPresenter {
	return ViewPresenter{w, u, e.NewErrorPresenter(w)}
}
func (p ViewPresenter) SuccessFindCollection(c clc.Collection) {
	MakeRow(p.RowURL, c).Render(p.Writer)
}

type EditPresenter struct {
	writer http.ResponseWriter
	b.RowURL
	task          string
	okMessageFunc func(update.Response) string
	htmx.ErrorPresenter
}

func NewEditPresenter(w http.ResponseWriter, u b.RowURL) EditPresenter {
	task := "Updating collection"
	okMessageFunc := func(r update.Response) string {
		return "Successfully updated collection"
	}
	return EditPresenter{w, u, task, okMessageFunc, htmx.NewErrorPresenter(task, w)}
}

func (p EditPresenter) SuccessUpdateCollection(r update.Response) {
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, p.okMessageFunc(r))
}

func (p EditPresenter) SuccessFindCollection(c clc.Collection) {
	b := bf.NewHTMXInlineFormBuilder(c.Name, len(listCollectionsFields), p.Url)
	b.AddTitle(fmt.Sprintf("Editing %v", c.Name))
	b.AddTextField("name", "Name", bf.WithRequired(), bf.WithDefault(c.Name))
	b.AddTextField("description", "Description", bf.WithDefault(c.Description))
	b.Render(p.writer)
}

type DeletePresenter struct {
	writer http.ResponseWriter
	b.RowURL
	task          string
	okMessageFunc func(update.Response) string
	htmx.ErrorPresenter
}

func NewDeletePresenter(w http.ResponseWriter, u b.RowURL) DeletePresenter {
	task := "Deleting collection"
	okMessageFunc := func(r update.Response) string {
		return "Successfully deleted collection"
	}
	return DeletePresenter{w, u, task, okMessageFunc, htmx.NewErrorPresenter(task, w)}
}

func (p DeletePresenter) SuccessFindCollection(c clc.Collection) {
	b.RenderConfirmDeleteRow(len(listCollectionsFields),
		c.Name, "collection", p.Url, p.writer)
}

func MakeRow(u b.RowURL, c clc.Collection) tb.Row {
	var groupName string
	if c.Group == nil {
		groupName = "n/a"
	} else {
		groupName = c.Group.Name
	}

	u.SetId(c.Name)
	actions := b.NewActionsPanelBuilder()
	actions.SetEdit(u.SetMode(b.ModeEdit).Url)
	actions.SetConfirmDelete(u.SetMode(b.ModeConfirmDelete).Url)

	row := tb.NewRow()
	row.AddCell(tb.NewCell(cmp.MakeTextLink(rt.MakeImagesURL(c.Name), c.Name)))
	row.AddCell(tb.NewCell(Text(c.Description)))
	row.AddCell(tb.NewCell(Text(groupName)))
	row.AddCell(tb.NewCell(Text(cmp.DateTimeToStr(c.CreatedAt))))
	row.AddCell(tb.NewCell(actions.Build()))
	return row

}
