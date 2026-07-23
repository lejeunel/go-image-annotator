package table

import (
	"io"
	"strings"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type TableBuilder struct {
	fields []string
	rows   []Row
}

func NewTableBuilder(fields []string) TableBuilder {
	return TableBuilder{fields: fields}
}
func (t TableBuilder) NumRows() int {
	return len(t.rows)
}

func (t *TableBuilder) Build() Node {
	if len(t.rows) == 0 {
		return Div(Class("ml-8"),
			Div(Class("pt-2 pb-8 text-lg italic"), Text("There is nothing here yet...")),
			Pre(Class("font-mono whitespace-pre text-sm leading-tight"), Raw(emptyAsciiIcon)),
		)
	}
	return Div(Class("overflow-hidden w-full overflow-x-auto rounded-radius border border-outline dark:border-outline-dark"),
		Table(Class("table-fixed w-full text-left text-sm text-on-surface dark:text-on-surface-dark"),
			TableHeader(t.fields),
			TableBody(t.rows),
		))
}

func (t *TableBuilder) AddRow(r Row) {
	t.rows = append(t.rows, r)
}

type Row struct {
	Cells []Cell
}

func NewRow() Row {
	return Row{}
}

func (r *Row) AddCell(c Cell) *Row {
	r.Cells = append(r.Cells, c)
	return r
}

func (r Row) Build() Node {
	return Tr(
		Class("even:bg-primary/5 dark:even:bg-primary-dark/10"),
		Map(r.Cells, func(c Cell) Node {
			return Td(
				Class(strings.Join([]string{"p-2", c.ExtraClass}, " ")),
				Attr(c.ExtraAttr),
				c.Content)
		}))

}

func (r Row) Render(w io.Writer) {
	r.Build().Render(w)
}

func TableHeader(fields []string) Node {
	return THead(Tr(Class("border-b border-outline bg-surface-alt text-sm text-on-surface-strong dark:border-outline-dark dark:bg-surface-dark-alt dark:text-on-surface-dark-strong"),
		Map(fields, func(f string) Node {
			return Th(Scope("col"), Class("p-2"), Text(f))
		})))
}

func TableBody(rows []Row) Node {
	return TBody(
		Class("divide-y divide-outline dark:divide-outline-dark"),
		Attr("hx-target", "closest tr"),
		Attr("hx-swap", "outerHTML"),
		Map(rows, func(r Row) Node {
			return r.Build()

		}),
	)
}
