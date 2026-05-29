package presenters

type Colorizer interface {
	Colorize(string) string
}

type CyclicColorizer struct {
	Palette  []string
	ColorMap map[string]string
}

func NewCyclicColorizer() CyclicColorizer {
	return CyclicColorizer{Palette: Palette, ColorMap: make(map[string]string)}
}

func (c CyclicColorizer) Colorize(key string) string {
	color, ok := c.ColorMap[key]
	if ok {
		return color
	}
	newColor := c.Palette[len(c.ColorMap)%(len(c.Palette)-1)]
	c.ColorMap[key] = newColor
	return newColor
}
