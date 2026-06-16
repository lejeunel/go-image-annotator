package annotator

import (
	"fmt"
	"time"

	html "github.com/lejeunel/go-image-annotator/adapters/web/html"
	"github.com/lejeunel/go-image-annotator/modules/annotator/view"
	. "maragu.dev/gomponents"
)

type ImageInfosView struct {
	result Node
}

func (p *ImageInfosView) Build(info view.ImageInfo) Node {
	s := html.SpecCard{}
	s.Fields = append(s.Fields, html.SpecFields{Name: "id", Value: info.Id},
		html.SpecFields{Name: "collection", Value: info.Collection},
		html.SpecFields{Name: "mimetype", Value: info.Specs.MIMEType},
		html.SpecFields{Name: "dimensions",
			Value: fmt.Sprintf("%vx%v", info.Specs.Width, info.Specs.Height)},
		html.SpecFields{Name: "ingested", Value: info.Specs.IngestedAt.Format(time.DateOnly)})
	return s.Render()

}
