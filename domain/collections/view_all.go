package collections

import (
	ctx "context"
	g "datahub/generic"
	"net/http"
	"net/url"

	gt "maragu.dev/gomponents"
	gh "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
)

type CollectionsListViewer struct {
	Viewer               *g.GenericViewer
	CollectionService    *Service
	PageSize             int
	PaginationWidgetSize int
}

func (v *CollectionsListViewer) List(collections []Collection) gt.Node {
	var rows []g.TableRow
	for _, c := range collections {
		link := gh.A(gh.Href("/collection/"+c.Id.String()),
			gh.Class(g.UrlClass),
			gt.Text(c.Name))

		rows = append(rows, g.TableRow{Values: []gt.Node{link, gt.Raw(c.Group), gt.Raw(c.CreatedAt.Format("2006-01-02 / 15:04")), gt.Raw(c.Id.String())}})
	}

	return g.MyTable([]string{"name", "group", "created_at", "id"},
		rows)
}

func (v *CollectionsListViewer) Render(ctx ctx.Context, page int64) gt.Node {
	collections, pagination, err := v.CollectionService.List(ctx,
		AlphabeticalOrdering,
		g.PaginationParams{Page: page,
			PageSize: 10})
	if err != nil {
		return v.Viewer.OopsPage(ctx, err.Error())
	}

	title := "Collections"

	baseURL, _ := url.Parse("/collections")
	paginatedViewer, err := g.MakePaginationViewer(*pagination, *baseURL, v.PageSize,
		v.PaginationWidgetSize)
	if err != nil {
		return v.Viewer.OopsPage(ctx, err.Error())
	}
	paginatedViewerPanel := gt.Group([]gt.Node{
		*paginatedViewer,
		v.List(collections),
		*paginatedViewer})

	body := gh.Div(gh.Class("px-2 md:px-4 lg:px8 py-2 md:py-4"), paginatedViewerPanel)
	return v.Viewer.BasePage(ctx, title, g.Head(), body)
}

func (v *CollectionsListViewer) Handler() http.HandlerFunc {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (gt.Node, error) {

		return v.Render(r.Context(), g.GetPage(*r)), nil
	})
}
