package builders

import (
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	html "github.com/lejeunel/go-image-annotator/adapters/web/html"
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
	tableBuilder   TableBuilder
	listURL        string
	pagination     s.Pagination
	creationButton *CreationButton
}

func NewPaginatedListBuilder(fields []string, listURL string, pagination s.Pagination) PaginatedListBuilder {
	tableBuilder := NewTableBuilder(fields)
	return PaginatedListBuilder{tableBuilder, listURL, pagination, nil}
}

func (b *PaginatedListBuilder) AddRow(nodes ...Node) *PaginatedListBuilder {
	b.tableBuilder.AddRow(nodes...)

	return b
}

func (b *PaginatedListBuilder) AddCreationButton(buttonLabel string, formGetEndpoint string, formDivId string) *PaginatedListBuilder {
	b.creationButton = &CreationButton{buttonLabel, formGetEndpoint, formDivId}
	return b
}
func (b *PaginatedListBuilder) Build() Node {
	paginator := html.MakePaginator(b.listURL, int(b.pagination.Page),
		int(b.pagination.TotalPages), b.tableBuilder.NumRows(), int(b.pagination.TotalRecords))

	var creationPanel Node
	if b.creationButton != nil {
		button := cmp.MakeHTMXCreateButton("Create new collection", b.creationButton.formGetEndpoint, b.creationButton.formDivId)
		formPlaceholder := Div(ID(b.creationButton.formDivId))
		creationPanel = Div(button, formPlaceholder)

	}
	return Div(creationPanel, Div(Class("py-2"), paginator), b.tableBuilder.Build())

}
