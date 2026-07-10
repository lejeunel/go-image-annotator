package builders

import (
	"net/url"

	ic "github.com/lejeunel/go-image-annotator/adapters/web/icons"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type Item struct {
	URL  url.URL
	Icon string
}

type ActionsPanelBuilder struct {
	Items []Item
}

func (p *ActionsPanelBuilder) SetEdit(url url.URL) *ActionsPanelBuilder {
	p.Items = append(p.Items, Item{Icon: ic.EditIcon, URL: url})
	return p
}

func (p *ActionsPanelBuilder) SetConfirmDelete(url url.URL) *ActionsPanelBuilder {
	q := url.Query()
	q.Set("confirm", "true")
	url.RawQuery = q.Encode()
	p.Items = append(p.Items, Item{Icon: ic.TrashIcon, URL: url})
	return p
}

func (p *ActionsPanelBuilder) Build() Node {
	res := []Node{}
	for _, a := range p.Items {
		var attr Node
		attr = Attr("hx-get", a.URL.String())
		res = append(res, A(
			Class("cursor-pointer"),
			attr,
			Raw(a.Icon)))
	}
	return Span(Class("inline-flex items-center gap-1"), Group(res))
}

func NewActionsPanelBuilder() ActionsPanelBuilder {
	return ActionsPanelBuilder{}
}
