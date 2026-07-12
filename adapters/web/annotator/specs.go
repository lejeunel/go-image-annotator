package annotator

import (
	"fmt"

	c "github.com/lejeunel/go-image-annotator/adapters/web/components"
	"github.com/lejeunel/go-image-annotator/modules/annotator/view"
	. "maragu.dev/gomponents"
)

type ImageInfosView struct {
	result Node
}

func (p *ImageInfosView) Build(info view.ImageInfo) Node {
	s := c.SpecCard{}
	s.Fields = append(s.Fields, c.SpecFields{Name: "id", Value: info.Id},
		c.SpecFields{Name: "collection", Value: info.Collection},
		c.SpecFields{Name: "mimetype", Value: info.Specs.MIMEType},
		c.SpecFields{Name: "dimensions",
			Value: fmt.Sprintf("%vx%v", info.Specs.Width, info.Specs.Height)},
		c.SpecFields{Name: "ingested", Value: c.DateTimeToStr(info.Specs.IngestedAt)})
	return s.Render()

}
