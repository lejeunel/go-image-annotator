package docs

import (
	"bytes"
	"embed"
	"errors"
	"io"
	"sort"
)

//go:embed pages
var docsFiles embed.FS

var ErrPageListEmpty = errors.New("no documentation pages provided")
var ErrParsing = errors.New("error reading documentation file")
var ErrMetaDataParsing = errors.New("error parsing meta-data from yaml section")

type Builder struct {
	parser Parser
	pages  []io.Reader
}

func NewBuilder(parser Parser) *Builder {
	return &Builder{parser: parser}
}

type Page struct {
	Title   string
	Weight  int
	Content bytes.Buffer
}

func (b *Builder) Build() ([]*Page, error) {
	if len(b.pages) == 0 {
		return nil, ErrPageListEmpty
	}

	var result []*Page
	for _, page := range b.pages {
		parsed, err := b.parser.Parse(page)
		if err != nil {
			return nil, err
		}

		result = append(result, parsed)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Weight < result[j].Weight
	})
	return result, nil
}

func (b *Builder) AddPage(page io.Reader) {
	b.pages = append(b.pages, page)
}
