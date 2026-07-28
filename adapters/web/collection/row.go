package collection

import (
	"io"
	"net/http"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	tb "github.com/lejeunel/go-image-annotator/adapters/web/builders/table"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	e "github.com/lejeunel/go-image-annotator/adapters/web/error"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	rt "github.com/lejeunel/go-image-annotator/routes"
	. "maragu.dev/gomponents"
)

type ViewPresenter struct {
	Writer io.Writer
	b.RowURL
	e.ErrorPresenter
}

func NewViewPresenter(w http.ResponseWriter, u b.RowURL) ViewPresenter {
	return ViewPresenter{w, u, e.NewErrorPresenter(w)}
}
func (p ViewPresenter) SuccessFindCollection(c clc.Collection) {
	MakeRow(p.RowURL, c).Render(p.Writer)
}

func MakeRow(u b.RowURL, c clc.Collection) tb.Row {
	var groupName string
	if c.Group == nil {
		groupName = "n/a"
	} else {
		groupName = c.Group.Name
	}

	u.SetId(c.Name)
	actions := b.NewActionsPanelBuilder()
	actions.SetEdit(u.SetMode(b.ModeEdit).Url)
	actions.SetConfirmDelete(u.SetMode(b.ModeConfirmDelete).Url)
	actions.SetClone(u.SetMode(b.ModeClone).Url)

	row := tb.NewRow()
	row.AddCell(tb.NewCell(cmp.MakeTextLink(rt.MakeImagesURL(c.Name), c.Name)))
	row.AddCell(tb.NewCell(Text(c.Description)))
	row.AddCell(tb.NewCell(Text(groupName)))
	row.AddCell(tb.NewCell(Text(cmp.DateTimeToStr(c.CreatedAt))))
	row.AddCell(tb.NewCell(actions.Build()))
	return row

}
