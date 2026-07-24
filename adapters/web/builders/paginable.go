package builders

import (
	"io"

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
	hasCreationButton bool
	creationButton    CreationButton
	PaginableTableBuilder
	PageBuilder
}

func NewPaginatedListBuilder(base PageBuilder, fields []string) PaginatedListBuilder {
	return PaginatedListBuilder{PaginableTableBuilder: NewPaginableTableBuilder(fields),
		PageBuilder: base}
}

func (b *PaginatedListBuilder) Render(w io.Writer) {
	paginator := cmp.MakePaginator(b.listURL, int(b.pagination.Page),
		int(b.pagination.TotalPages), b.TableBuilder.NumRows(), int(b.pagination.TotalRecords))

	var creationPanel Node
	if b.hasCreationButton {
		button := cmp.MakeHTMXCreateButton(b.creationButton.label, b.creationButton.formGetEndpoint, b.creationButton.formDivId)
		formPlaceholder := Div(ID(b.creationButton.formDivId))
		creationPanel = Div(button, formPlaceholder)
	}
	content := Div(
		creationPanel,
		Div(Class("py-2"), paginator),
		b.PaginableTableBuilder.Build())

	b.PageBuilder.SetContent(content)
	b.PageBuilder.Render(w)
}
func (b *PaginatedListBuilder) AddCreationButton(buttonLabel string, formEndpoint string, formDivId string) *PaginatedListBuilder {
	b.creationButton = CreationButton{buttonLabel, formEndpoint, formDivId}
	b.hasCreationButton = true
	return b
}
