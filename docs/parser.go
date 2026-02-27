package docs

import (
	"bytes"
	"fmt"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"io"
)

type Parser interface {
	Parse(io.Reader) (*Page, error)
}

type GoldMarkParser struct {
	converter goldmark.Markdown
}

func NewGoldMarkParser() *GoldMarkParser {
	return &GoldMarkParser{converter: goldmark.New(
		goldmark.WithExtensions(meta.Meta),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	)}
}

func (p *GoldMarkParser) Parse(input io.Reader) (*Page, error) {
	source, err := io.ReadAll(input)
	if err != nil {
		return nil, ErrParsing
	}
	var destination bytes.Buffer
	context := parser.NewContext()
	if err := p.converter.Convert([]byte(source), &destination, parser.WithContext(context)); err != nil {
		panic(err)
	}
	page, err := p.initMetaData(context)
	if err != nil {
		return nil, err
	}

	page.Content = destination
	return page, nil

}

func (p *GoldMarkParser) initMetaData(ctx parser.Context) (*Page, error) {
	metaData := meta.Get(ctx)
	title, ok := metaData["Title"].(string)
	if !ok {
		return nil, fmt.Errorf("%v: %w", "missing or invalid Title field in yaml section", ErrMetaDataParsing)
	}

	weight, ok := metaData["Weight"].(int)
	if !ok {
		return nil, fmt.Errorf("%v: %w", "missing or invalid Weight field in yaml section", ErrMetaDataParsing)
	}
	return &Page{Title: title, Weight: weight}, nil
}
