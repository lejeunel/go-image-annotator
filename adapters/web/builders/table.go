package builders

import (
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
	return Div(Class("overflow-hidden w-full overflow-x-auto rounded-radius border border-outline dark:border-outline-dark"),
		Table(Class("w-full text-left text-sm text-on-surface dark:text-on-surface-dark"),
			TableHeader(t.fields),
			TableBody(t.rows),
		))
}

func (t *TableBuilder) AddRow(nodes ...Node) {
	t.rows = append(t.rows, Row{nodes})
}

type Row struct {
	Values []Node
}

func (r Row) Render() Node {
	return Tr(
		Class("even:bg-primary/5 dark:even:bg-primary-dark/10"),
		Map(r.Values, func(node Node) Node {
			return Td(Class("p-2"),
				node)
		}))

}

func TableHeader(fields []string) Node {
	return THead(Tr(Class("border-b border-outline bg-surface-alt text-sm text-on-surface-strong dark:border-outline-dark dark:bg-surface-dark-alt dark:text-on-surface-dark-strong"),
		Map(fields, func(f string) Node {
			return Th(Scope("col"), Class("p-2"), Text(f))
		})))
}

func TableBody(rows []Row) Node {
	return TBody(Class("divide-y divide-outline dark:divide-outline-dark"),
		Map(rows, func(r Row) Node {
			return r.Render()

		}),
	)
}
