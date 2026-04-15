package html

import (
	gp "maragu.dev/gomponents"
	gh "maragu.dev/gomponents/html"
)

type SpecTable struct {
	Rows []SpecTableRow
}

func (t *SpecTable) Render() gp.Node {
	// return gh.Div(gh.Class("inline-block overflow-hidden overflow-x-auto rounded-radius border border-outline dark:border-outline-dark"),
	return gh.Div(gh.Class("rounded-radius border border-outline dark:border-outline-dark"),
		gh.Table(gh.Class("text-left text-sm text-on-surface dark:text-on-surface-dark"),
			gp.Map(t.Rows, func(r SpecTableRow) gp.Node {
				return r.Render()
			}),
		))
}

type SpecTableRow struct {
	Name  string
	Value string
}

func (r SpecTableRow) Render() gp.Node {
	return gh.Tr(gh.Td(gh.Class("py-2 px-2 font-bold"), gp.Text(r.Name)),
		gh.Td(gh.Class("py-2 px-2"), gp.Text(r.Value)))

}
