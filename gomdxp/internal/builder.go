package docs

import (
	"io"
	"sort"
)

type Builder struct {
	parser Parser
	inputs []io.Reader
	pages  []*Page
}

func NewBuilder(parser Parser) *Builder {
	return &Builder{parser: parser}
}

func (b *Builder) Parse() error {
	if len(b.inputs) == 0 {
		return ErrPageListEmpty
	}

	for _, page := range b.inputs {
		parsed, err := b.parser.Parse(page)
		if err != nil {
			return err
		}
		b.pages = append(b.pages, parsed)
	}

	sort.Slice(b.pages, func(i, j int) bool {
		return b.pages[i].Weight < b.pages[j].Weight
	})
	return nil
}

func (b *Builder) Build() []*Page {
	return b.pages
}

func (b *Builder) AddPage(page io.Reader) {
	b.inputs = append(b.inputs, page)
}
