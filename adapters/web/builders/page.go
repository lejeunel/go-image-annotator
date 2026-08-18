package builders

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"maps"

	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	g "github.com/lejeunel/go-image-annotator/globals"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"github.com/yuin/goldmark"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

//go:embed scripts/detect_os.js
var detectOs string

//go:embed templates/query_modal.html
var queryModal string

type QueryModalData struct {
	SubmitURL            string
	FilterQueryArgName   string
	OrderingQueryArgName string
}

type PageBuilder struct {
	APIPath             string
	RepoURL             string
	DocsURL             string
	Version             g.Info
	ActivePage          cmp.ActivePage
	User                *u.User
	SidebarTitle        string
	SidebarEntries      map[string]cmp.SidebarEntry
	SidebarEntriesOrder []string
	markdownPreamble    string
	postamble           string
	content             Node
	Title               string
	BasePageBuilder
}

func NewPageBuilder(base BasePageBuilder, version g.Info) PageBuilder {
	pb := PageBuilder{
		BasePageBuilder: base, APIPath: rt.APIRootUrl, RepoURL: g.RepoURL, DocsURL: g.DocsURL,
		Version: version, SidebarEntries: make(map[string]cmp.SidebarEntry),
	}
	pb.AddScripts(Script(Raw(detectOs)))
	return pb
}

func (b *PageBuilder) SetHTMLTitle(title string) *PageBuilder {
	b.BasePageBuilder.SetHTMLTitle(title)
	return b
}

func (b *PageBuilder) SetTitle(title string) *PageBuilder {
	b.Title = title
	return b
}

func (b *PageBuilder) SetActiveSection(a cmp.ActivePage) *PageBuilder {
	b.ActivePage = a
	return b
}

func (b *PageBuilder) SetUserIdentity(ctx context.Context) *PageBuilder {
	id := u.IdentityFromContext(ctx)
	b.User = id
	return b
}

func (b *PageBuilder) AddSidebarTitle(title string) *PageBuilder {
	b.SidebarTitle = title
	return b
}

func (b *PageBuilder) ActivateSidebarEntry(name string) *PageBuilder {
	b.SidebarEntries = maps.Clone(b.SidebarEntries)
	for k, v := range b.SidebarEntries {
		v.IsActive = false
		if k == name {
			v.IsActive = true
		}
		b.SidebarEntries[k] = v

	}
	return b
}

func (b *PageBuilder) AddSidebarEntry(name, icon, url string, isActive bool) *PageBuilder {
	b.SidebarEntries = maps.Clone(b.SidebarEntries)
	b.SidebarEntries[name] = cmp.SidebarEntry{Label: name, Icon: icon, Url: url, IsActive: isActive}
	b.SidebarEntriesOrder = append(b.SidebarEntriesOrder, name)
	return b
}

func (b *PageBuilder) markdownToHTML(data string) string {
	md := goldmark.New()
	var buf bytes.Buffer
	if err := md.Convert([]byte(data), &buf); err != nil {
		panic(err)
	}
	return buf.String()
}

func (b *PageBuilder) AddMarkdownPreamble(md string) *PageBuilder {
	b.markdownPreamble = b.markdownToHTML(md)
	return b
}

func (b *PageBuilder) AddMarkdownPostamble(md string) *PageBuilder {
	b.postamble = b.markdownToHTML(md)
	return b
}

func (b *PageBuilder) SetContent(content Node) *PageBuilder {
	b.content = content
	return b
}

func (b *PageBuilder) Render(w io.Writer) {
	if b.User == nil {
		b.BasePageBuilder.SetError(fmt.Errorf("current user has not been set"))
		b.BasePageBuilder.Render(w)
		return
	}

	var content Node

	if b.Title != "" {
		content = Div(
			content,
			Div(
				Class(
					"text-2xl/7 font-bold text-white sm:truncate sm:text-3xl sm:tracking-tight font-roboto",
				),
				Text(b.Title),
			),
		)
	}

	if b.markdownPreamble != "" {
		content = Div(
			content,
			Div(
				Class("flex flex-col w-full"),
				Article(Class("prose dark:prose-invert max-w-none mb-4"), Raw(b.markdownPreamble)),
			),
		)
	}

	// content = Div(Class("flex flex-col w-full"), content)
	content = Div(Class("flex-1 min-w-0"), content, b.content)

	if b.postamble != "" {
		content = Div(content, Div(Class("flex flex-col w-150 mt-4"), cmp.Separator,
			Article(Class("prose dark:prose-invert max-w-none"), Raw(b.postamble))))
	}

	if len(b.SidebarEntries) > 0 {
		var bufSidebar bytes.Buffer
		sidebar := cmp.NewSidebar(b.SidebarTitle)
		for _, n := range b.SidebarEntriesOrder {
			e, _ := b.SidebarEntries[n]
			sidebar.AddEntry(e.Label, e.Icon, e.Url, e.IsActive)
		}
		sidebar.Render(&bufSidebar)
		content = Div(
			Class("relative flex w-full flex-col"),
			Nav(
				Attr("x-cloak"),
				Class(
					`fixed left-0 top-14 z-20 flex h-svh w-60 shrink-0 flex-col border-r border-outline bg-surface-alt p-4 transition-transform duration-300
                      dark:border-outline-dark dark:bg-surface-dark-alt`,
				),
				Raw(bufSidebar.String()),
			),
			Div(Class("ml-60 px-4 py-18"), content),
		)
	} else {
		content = Div(Class("flex overflow-hidden w-full px-4 py-18"), content)
	}

	queryModalTemplate, err := template.New("").Parse(queryModal)
	if err != nil {
		Text(err.Error()).Render(w)
		return
	}
	var queryBuf bytes.Buffer
	if err := queryModalTemplate.ExecuteTemplate(
		&queryBuf,
		"query",
		QueryModalData{
			SubmitURL:            rt.SliceUrl,
			FilterQueryArgName:   rt.FilterQueryArgName,
			OrderingQueryArgName: rt.OrderingQueryArgName,
		}); err != nil {
		Text(err.Error()).Render(w)
		return
	}

	b.BasePageBuilder.SetFrameContent(
		Div(
			Attr(`x-data="{ showSearch: false}"`),
			Attr(`@keydown.cmd.k.window.prevent="showSearch = !showSearch, $dispatch('searchModalOpened')"`),
			Attr(`@keydown.ctrl.k.window.prevent="showSearch = !showSearch, $dispatch('searchModalOpened')"`),
			Group(
				[]Node{
					cmp.MakeNavBar(
						b.ActivePage,
						b.RepoURL,
						b.DocsURL,
						b.APIPath,
						*b.User,
						rt.DashboardUrl,
					),
					content,
					cmp.MakeFooter(b.Version),
				},
			),
			Raw(queryBuf.String()),
		),
	)
	b.BasePageBuilder.Render(w)
}
