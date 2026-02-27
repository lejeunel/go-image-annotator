package annotator

// https://www.color-hex.com/color-palette/10221
var Palette = []string{"#ff71ce", "#01cdfe", "#05ffa1",
	"#b967ff", "#fffb96"}

type Colorizer struct {
	Colors []string
}

func (c *Colorizer) Colorize(boxes []*BoundingBox) {
	for i, b := range boxes {
		b.Color = c.Colors[i%(len(c.Colors)-1)]
	}
}
