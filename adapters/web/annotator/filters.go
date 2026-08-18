package annotator

import (
	c "github.com/lejeunel/go-image-annotator/adapters/web/components"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	. "maragu.dev/gomponents"
)

type QueryView struct {
	result Node
}

func (p *QueryView) Build(f im.FilterStr, o im.OrderStr) Node {
	s := c.SpecCard{}
	s.Fields = append(s.Fields,
		c.SpecFields{Name: "filters", Value: f},
		c.SpecFields{Name: "ordering", Value: o})
	return s.Render()
}
