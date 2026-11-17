package components

import (
	gp "maragu.dev/gomponents"
	gh "maragu.dev/gomponents/html"
)

type SpecCard struct {
	Button *gp.Node
	Fields []SpecFields
}

func (t *SpecCard) Render() gp.Node {
	return gh.Div(
		gh.Class("flex items-start rounded-radius border border-outline dark:border-outline-dark"),
		gh.Table(
			gh.Class("w-full flex-1 text-left text-sm text-on-surface dark:text-on-surface-dark"),
			gp.Map(t.Fields, func(r SpecFields) gp.Node {
				return r.Render()
			}),
		),
		gp.Iff(t.Button != nil, func() gp.Node {
			return gh.Div(gh.Class("shrink-0 p-2 leading-none"), *t.Button)
		}),
	)
}

type SpecFields struct {
	Name  string
	Value string
}

func (r SpecFields) Render() gp.Node {
	return gh.Tr(gh.Td(gh.Class("w-px whitespace-nowrap py-2 px-2 font-bold"), gp.Text(r.Name)),
		gh.Td(gh.Class("py-2 px-2"), gp.Text(r.Value)))
}
