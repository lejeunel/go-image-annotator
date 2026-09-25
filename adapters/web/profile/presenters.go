package profile

import (
	_ "embed"
	"io"
	"net/http"
	"strings"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	tb "github.com/lejeunel/go-image-annotator/adapters/web/builders/table"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	e "github.com/lejeunel/go-image-annotator/adapters/web/error"
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"github.com/lejeunel/go-image-annotator/use-cases/profile/list"
	. "maragu.dev/gomponents"
)

var listProfilesFields = []string{"name", "description", "labels", "actions"}

type ListPresenter struct {
	b.PaginatedListBuilder
	b.RowURL
	io.Writer
	e.ErrorPresenter
}

func NewListPresenter(w http.ResponseWriter, p b.PageBuilder, u b.RowURL) ListPresenter {
	p.SetTitle("Profiles").SetHTMLTitle("Profiles").SetActiveSection(cmp.ProfilesPageActive)
	pb := b.NewPaginatedListBuilder(p, listProfilesFields)
	return ListPresenter{pb, u, w, e.NewErrorPresenter(w, p)}
}

func (p ListPresenter) SuccessListProfiles(r list.Response) {
	p.SetPagination(r.Pagination, rt.ProfilesUrl)
	for _, l := range r.Profiles {
		row := MakeRow(p.RowURL, l)
		p.AddRow(row)
	}

	p.AddCreationButton("Create", CreateProfileFormUrl, createProfileTargetDiv)
	p.PaginatedListBuilder.AddMarkdownPreamble(preamble)
	p.Render(p.Writer)
}

type ViewPresenter struct {
	io.Writer
	b.RowURL
	e.ErrorPresenter
}

func NewViewPresenter(w http.ResponseWriter, p b.PageBuilder, u b.RowURL) ViewPresenter {
	return ViewPresenter{w, u, e.NewErrorPresenter(w, p)}
}

func (p ViewPresenter) SuccessFindProfile(l pr.Profile) {
	MakeRow(p.RowURL, l).Render(p.Writer)
}

type EditPresenter struct {
	io.Writer
	b.RowURL
	e.ErrorPresenter
	availableLabels []string
}

type DeletePresenter struct {
	io.Writer
	b.RowURL
	e.ErrorPresenter
}

func NewDeletePresenter(w http.ResponseWriter, p b.PageBuilder, u b.RowURL) DeletePresenter {
	return DeletePresenter{w, u, e.NewErrorPresenter(w, p)}
}

func (p DeletePresenter) SuccessFindProfile(l pr.Profile) {
	b.RenderConfirmDeleteRow(len(listProfilesFields),
		l.Name, "profile", p.Url, p.Writer)
}

func MakeRow(u b.RowURL, l pr.Profile) tb.Row {
	u.SetId(l.Name)
	actions := b.NewActionsPanelBuilder()
	actions.SetEdit(u.SetMode(b.ModeEdit).Url)
	actions.SetConfirmDelete(u.SetMode(b.ModeConfirmDelete).Url)
	row := tb.NewRow()
	row.AddCell(tb.NewCell(Text(l.Name)))
	row.AddCell(tb.NewCell(Text(l.Description)))
	row.AddCell(tb.NewCell(Text(strings.Join(l.Labels, " / "))))
	row.AddCell(tb.NewCell(actions.Build()))

	return row
}
