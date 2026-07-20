package builders

import (
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type CreationButton struct {
	label           string
	formGetEndpoint string
	formDivId       string
}

type PaginatedListBuilder struct {
	creationButton *CreationButton
	PaginableTableBuilder
	PageBuilder
}

func NewPaginatedListBuilder(base PageBuilder, fields []string) PaginatedListBuilder {
	return PaginatedListBuilder{PaginableTableBuilder: NewPaginableTableBuilder(fields),
		PageBuilder: base}
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
	content := Div(creationPanel, Div(Class("py-2"), paginator), b.PaginableTableBuilder.Build())
	b.SetContent(content)
	return b

}
