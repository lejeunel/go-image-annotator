package html

import (
	ic "github.com/lejeunel/go-image-annotator/adapters/web/icons"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type Item struct {
	URL  string
	Icon string
}

type ActionsPanel struct {
	Items []Item
}

func (p *ActionsPanel) SetEdit(url string) *ActionsPanel {
	p.Items = append(p.Items, Item{Icon: ic.EditIcon, URL: url})
	return p
}

func (p *ActionsPanel) SetDelete(url string) *ActionsPanel {
	p.Items = append(p.Items, Item{Icon: ic.TrashIcon, URL: url})
	return p
}

func (p *ActionsPanel) Build() Node {
	res := []Node{}
	for _, a := range p.Items {
		res = append(res, A(
			Class("cursor-pointer"),
			Attr("hx-get", a.URL),
			Raw(a.Icon)))
	}
	return Span(Class("inline-flex items-center gap-1"), Group(res))
}

func NewActionsPanel() ActionsPanel {
	return ActionsPanel{}
}
