package annotator

import (
	"fmt"
	"time"

	html "github.com/lejeunel/go-image-annotator/adapters/web/html"
	"github.com/lejeunel/go-image-annotator/app/annotator/view"
	. "maragu.dev/gomponents"
)

type ImageInfosView struct {
	result Node
}

func (p *ImageInfosView) Build(info view.ImageInfo) Node {
	table := html.SpecTable{}
	table.Rows = append(table.Rows, html.SpecTableRow{Name: "id", Value: ShortenUUID(info.Id)})
	table.Rows = append(table.Rows, html.SpecTableRow{Name: "collection", Value: info.Collection})
	table.Rows = append(table.Rows, html.SpecTableRow{Name: "mimetype", Value: info.Specs.MIMEType})
	table.Rows = append(table.Rows, html.SpecTableRow{Name: "dimensions",
		Value: fmt.Sprintf("%vx%v", info.Specs.Width, info.Specs.Height)})
	table.Rows = append(table.Rows, html.SpecTableRow{Name: "ingested", Value: info.Specs.IngestedAt.Format(time.DateOnly)})
	return table.Render()

}
