package builders

import (
	"net/url"

	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
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

func (p *ActionsPanelBuilder) SetClone(url url.URL) *ActionsPanelBuilder {
	p.Items = append(p.Items, Item{Icon: ic.Copy, URL: url, Tooltip: "clone"})
	return p
}

func (p *ActionsPanelBuilder) SetExpand(url url.URL) *ActionsPanelBuilder {
	p.Items = append(p.Items, Item{Icon: ic.Expand, URL: url, Tooltip: "expand"})
	return p
}

func (p *ActionsPanelBuilder) Build() Node {
	res := []Node{}
	for _, a := range p.Items {
		res = append(res, cmp.MakeIconizedButton(a.Icon, a.Tooltip, Attr("hx-get", a.URL.String())))
	}
	return Span(Class("inline-flex items-center gap-1"), Group(res))
}

func NewActionsPanelBuilder() ActionsPanelBuilder {
	return ActionsPanelBuilder{}
}
