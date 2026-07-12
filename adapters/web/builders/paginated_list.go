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
	tableBuilder   tb.TableBuilder
	listURL        string
	pagination     s.Pagination
	creationButton *CreationButton
}

func NewPaginatedListBuilder(fields []string, listURL string, pagination s.Pagination) PaginatedListBuilder {
	tableBuilder := tb.NewTableBuilder(fields)
	return PaginatedListBuilder{tableBuilder, listURL, pagination, nil}
}

func (b *PaginatedListBuilder) AddRow(r tb.Row) *PaginatedListBuilder {
	b.tableBuilder.AddRow(r)

	return b
}

func (b *PaginatedListBuilder) AddCreationButton(buttonLabel string, formEndpoint string, formDivId string) *PaginatedListBuilder {
	b.creationButton = &CreationButton{buttonLabel, formEndpoint, formDivId}
	return b
}
func (b *PaginatedListBuilder) Build() Node {
	paginator := cmp.MakePaginator(b.listURL, int(b.pagination.Page),
		int(b.pagination.TotalPages), b.tableBuilder.NumRows(), int(b.pagination.TotalRecords))

	var creationPanel Node
	if b.creationButton != nil {
		button := cmp.MakeHTMXCreateButton("Create new collection", b.creationButton.formGetEndpoint, b.creationButton.formDivId)
		formPlaceholder := Div(ID(b.creationButton.formDivId))
		creationPanel = Div(button, formPlaceholder)

	}
	return Div(creationPanel, Div(Class("py-2"), paginator), b.tableBuilder.Build())

}
