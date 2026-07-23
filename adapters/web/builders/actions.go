package builders

import (
	"net/url"

	ic "github.com/lejeunel/go-image-annotator/adapters/web/icons"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type Item struct {
	URL     url.URL
	Icon    string
	Tooltip string
}

type ActionsPanelBuilder struct {
	Items []Item
}

func (p *ActionsPanelBuilder) SetEdit(url url.URL) *ActionsPanelBuilder {
	p.Items = append(p.Items, Item{Icon: ic.Edit, URL: url, Tooltip: "edit"})
	return p
}

func (p *ActionsPanelBuilder) SetConfirmDelete(url url.URL) *ActionsPanelBuilder {
	p.Items = append(p.Items, Item{Icon: ic.Trash, URL: url, Tooltip: "delete"})
	return p
}

func (p *ActionsPanelBuilder) Build() Node {
	res := []Node{}
	for _, a := range p.Items {
		var attr Node
		attr = Attr("hx-get", a.URL.String())
		res = append(res,
			Div(
				Class("relative w-fit"),
				Button(
					Type("button"),
					Class("peer rounded-radius bg-surface-alt border border-surface-alt px-1 py-1 font-medium tracking-wide text-on-surface focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary dark:bg-surface-dark-alt dark:border-surface-dark-alt dark:text-on-surface-dark dark:focus-visible:outline-primary-dark cursor-pointer"),
					attr,
					Raw(a.Icon)),
				Div(
					Class("absolute -top-9 left-1/2 -translate-x-1/2 z-10 whitespace-nowrap rounded-sm bg-surface-dark px-2 py-1 text-center text-sm text-on-surface-dark-strong opacity-0 transition-all ease-out peer-hover:opacity-100 peer-focus:opacity-100 pointer-events-none dark:bg-surface dark:text-on-surface-strong"),
					Role("tooltip"),
					Text(a.Tooltip))))
	}
	return Span(Class("inline-flex items-center gap-1"), Group(res))
}

func NewActionsPanelBuilder() ActionsPanelBuilder {
	return ActionsPanelBuilder{}
}
