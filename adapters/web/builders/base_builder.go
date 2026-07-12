package builders

import (
	"bytes"
	"io"

	_ "embed"
	"html/template"

	. "maragu.dev/gomponents"
)

//go:embed templates/base.html
var BaseHTML string

type BasePageBuilder struct {
	Title     string
	pane      *Node
	scripts   []Node
	bodyExtra []string
	Error     error
	Content   Node
}

func (b *BasePageBuilder) AddScripts(scripts ...Node) *BasePageBuilder {
	for _, s := range scripts {
		b.scripts = append(b.scripts, s)
	}
	return b
}
func (b *BasePageBuilder) SetTitle(title string) *BasePageBuilder {
	b.Title = title
	return b
}
func (b *BasePageBuilder) SetContent(c Node) *BasePageBuilder {
	b.Content = c
	return b
}
func (b *BasePageBuilder) SetError(err error) *BasePageBuilder {
	b.Error = err
	return b
}

type BaseData struct {
	Title   string
	Scripts template.HTML
	Content template.HTML
}

func (b *BasePageBuilder) Render(w io.Writer) {
	if b.Error != nil {
		Text(b.Error.Error()).Render(w)
		return
	}
	baseTemplate, err := template.New("").Parse(BaseHTML)
	if err != nil {
		Text(err.Error()).Render(w)
		return
	}

	var content bytes.Buffer
	if err := b.Content.Render(&content); err != nil {
		Text(err.Error()).Render(w)
	}
	var scripts bytes.Buffer
	if err := Group(b.scripts).Render(&scripts); err != nil {
		Text(err.Error()).Render(w)
		return
	}
	if err := baseTemplate.ExecuteTemplate(
		w,
		"base",
		BaseData{Title: b.Title,
			Content: template.HTML(content.String()),
			Scripts: template.HTML(scripts.String())}); err != nil {
		Text(err.Error()).Render(w)
	}

}

func NewBasePageBuilder() BasePageBuilder {
	return BasePageBuilder{}
}
