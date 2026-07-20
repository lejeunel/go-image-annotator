package builders

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"

	. "maragu.dev/gomponents"
)

type AdminPageBuilder struct {
	PageBuilder
}

//go:embed templates/admin.html
var AdminHTML string

func (b *AdminPageBuilder) Build() *AdminPageBuilder {
	content := Text("Hello")
	var buf bytes.Buffer
	if err := content.Render(&buf); err != nil {
		b.SetError(fmt.Errorf("rendering gomponents content: %w", err))
		return b
	}
	data := struct {
		Content template.HTML
	}{
		Content: template.HTML(buf.String()),
	}

	baseTemplate, _ := template.New("").Parse(AdminHTML)

	var out bytes.Buffer
	if err := baseTemplate.ExecuteTemplate(&out, "admin", data); err != nil {
		b.SetError(fmt.Errorf("executing admin template: %w", err))
		return b
	}

	b.SetContent(Raw(out.String()))
	return b
}

func NewAdminPageBuilder(b PageBuilder) AdminPageBuilder {
	return AdminPageBuilder{PageBuilder: b}
}
