package annotator

import (
	a "github.com/lejeunel/go-image-annotator-v2/application/annotator"
	html "github.com/lejeunel/go-image-annotator-v2/shared/html"
	. "maragu.dev/gomponents"
)

type ImageInfosView struct {
	result Node
}

func (p *ImageInfosView) Render(info a.ImageInfo) Node {
	table := html.SpecTable{}
	table.Rows = append(table.Rows, html.SpecTableRow{Name: "id", Value: info.Id})
	table.Rows = append(table.Rows, html.SpecTableRow{Name: "collection", Value: info.Collection})
	return table.Render()

}
