package html

import (
	gp "maragu.dev/gomponents"
	gh "maragu.dev/gomponents/html"
)

type PaginationTable struct {
	Fields []string
	Rows   []PaginationTableRow
}

func (t *PaginationTable) Build() gp.Node {
	return gh.Div(gh.Class("overflow-hidden w-full overflow-x-auto rounded-radius border border-outline dark:border-outline-dark"),
		gh.Table(gh.Class("w-full text-left text-sm text-on-surface dark:text-on-surface-dark"),
			TableHeader(t.Fields),
			TableBody(t.Rows),
		))
}

type PaginationTableRow struct {
	Values []gp.Node
}

func (r PaginationTableRow) Render() gp.Node {
	return gh.Tr(
		gh.Class("even:bg-primary/5 dark:even:bg-primary-dark/10"),
		gp.Map(r.Values, func(node gp.Node) gp.Node {
			return gh.Td(gh.Class("p-4"),
				node)
		}))

}

func TableHeader(fields []string) gp.Node {
	return gh.THead(gh.Tr(gh.Class("border-b border-outline bg-surface-alt text-sm text-on-surface-strong dark:border-outline-dark dark:bg-surface-dark-alt dark:text-on-surface-dark-strong"),
		gp.Map(fields, func(f string) gp.Node {
			return gh.Th(gh.Scope("col"), gh.Class("p-4"), gp.Text(f))
		})))
}

func TableBody(rows []PaginationTableRow) gp.Node {
	return gh.TBody(gh.Class("divide-y divide-outline dark:divide-outline-dark"),
		gp.Map(rows, func(r PaginationTableRow) gp.Node {
			return r.Render()

		}),
	)
}
