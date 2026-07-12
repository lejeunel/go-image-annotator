package builders

import (
	tb "github.com/lejeunel/go-image-annotator/adapters/web/builders/table"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	s "github.com/lejeunel/go-image-annotator/shared/pagination"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type CreationButton struct {
	label           string
	formGetEndpoint string
	formDivId       string
}

type PaginatedListBuilder struct {
	listURL        string
	pagination     s.Pagination
	creationButton *CreationButton
	tb.TableBuilder
	PageBuilder
}

func NewPaginatedListBuilder(base PageBuilder, fields []string) PaginatedListBuilder {
	tableBuilder := tb.NewTableBuilder(fields)
	return PaginatedListBuilder{TableBuilder: tableBuilder, PageBuilder: base}
}
func (b *PaginatedListBuilder) SetPagination(pg s.Pagination, url string) *PaginatedListBuilder {
	b.pagination = pg
	b.listURL = url
	return b
}

func (b *PaginatedListBuilder) AddRow(r tb.Row) *PaginatedListBuilder {
	b.TableBuilder.AddRow(r)
	return b
}

func (b *PaginatedListBuilder) AddCreationButton(buttonLabel string, formEndpoint string, formDivId string) *PaginatedListBuilder {
	b.creationButton = &CreationButton{buttonLabel, formEndpoint, formDivId}
	return b
}
func (b *PaginatedListBuilder) Build() *PaginatedListBuilder {
	paginator := cmp.MakePaginator(b.listURL, int(b.pagination.Page),
		int(b.pagination.TotalPages), b.TableBuilder.NumRows(), int(b.pagination.TotalRecords))

	var creationPanel Node
	if b.creationButton != nil {
		button := cmp.MakeHTMXCreateButton(b.creationButton.label, b.creationButton.formGetEndpoint, b.creationButton.formDivId)
		formPlaceholder := Div(ID(b.creationButton.formDivId))
		creationPanel = Div(button, formPlaceholder)

	}
	content := Div(creationPanel, Div(Class("py-2"), paginator), b.TableBuilder.Build())
	b.SetContent(content)
	return b

}
