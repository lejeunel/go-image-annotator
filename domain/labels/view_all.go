package labels

import (
	"context"
	g "datahub/generic"
	gp "maragu.dev/gomponents"
	gh "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
	"net/http"
	"net/url"
)

type LabelListViewer struct {
	Viewer               *g.GenericViewer
	LabelsService        *Service
	PageSize             int
	PaginationWidgetSize int
}

func (v LabelListViewer) List(labels []Label) gp.Node {
	var rows []g.TableRow
	for _, l := range labels {

		link := gh.A(gh.Href("/label/"+l.Id.String()),
			gh.Class(g.UrlClass),
			gp.Text(l.Name))
		rows = append(rows, g.TableRow{[]gp.Node{
			link,
			gp.Raw(l.Description),
			gp.Raw(""),
			gp.Raw(l.CreatedAt.Format("2006-01-02 / 15:04")),
			gp.Raw(l.Id.String()),
		}})
	}

	return g.MyTable([]string{"name", "description", "parent_id", "created_at", "id"},
		rows)

}

func (v *LabelListViewer) Render(ctx context.Context, page int64) gp.Node {
	labels, pagination, err := v.LabelsService.List(ctx,
		g.OrderingArg{Field: "name"},
		g.PaginationParams{Page: page,
			PageSize: v.PageSize})
	if err != nil {
		return v.Viewer.OopsPage(ctx, err.Error())
	}

	baseURL, _ := url.Parse("/labels")
	paginationWidget, err := g.MakePaginationViewer(*pagination, *baseURL, v.PageSize,
		v.PaginationWidgetSize)
	if err != nil {
		return v.Viewer.OopsPage(ctx, err.Error())
	}

	paginatedViewerPanel := gp.Group([]gp.Node{
		*paginationWidget,
		v.List(labels),
		*paginationWidget})

	body := gh.Div(gh.Class("px-2 md:px-4 lg:px8 py-2 md:py-4"), paginatedViewerPanel)

	return v.Viewer.BasePage(ctx, "Labels", g.Head(), body)
}

func (v *LabelListViewer) Handler() http.HandlerFunc {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (gp.Node, error) {
		return v.Render(r.Context(), g.GetPage(*r)), nil
	})
}
