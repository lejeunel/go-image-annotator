package builders

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"

	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	. "maragu.dev/gomponents"
)

//go:embed templates/sidebar.html
var SidebarHTML string

type SidebarItem struct {
	URL    string
	Label  string
	Icon   template.HTML
	Active bool
}

type SidebarPageData struct {
	Content template.HTML
	Sidebar template.HTML
}

type SideBarPageBuilder struct {
	PageBuilder
	Sidebar cmp.Sidebar
	content Node
}

func (b *SideBarPageBuilder) SetContent(c Node) *SideBarPageBuilder {
	b.content = c
	return b
}
func (b *SideBarPageBuilder) SetActiveSidebarItem(name string) *SideBarPageBuilder {
	b.Sidebar.Activate(name)
	return b
}
func (b *SideBarPageBuilder) Build() *SideBarPageBuilder {
	var contentBuf bytes.Buffer
	if b.content != nil {
		b.content.Render(&contentBuf)
	}
	var sidebarBuf bytes.Buffer
	b.Sidebar.Render(&sidebarBuf)
	baseTemplate, _ := template.New("").Parse(SidebarHTML)

	var out bytes.Buffer
	if err := baseTemplate.ExecuteTemplate(&out, "admin",
		SidebarPageData{
			template.HTML(contentBuf.String()),
			template.HTML(sidebarBuf.String())}); err != nil {
		b.SetError(fmt.Errorf("executing admin template: %w", err))
		return b
	}

	b.PageBuilder.SetContent(Raw(out.String()))
	return b
}
