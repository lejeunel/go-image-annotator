package locations

import (
	"context"
	g "datahub/generic"
	"net/http"
	"net/url"

	gt "maragu.dev/gomponents"
	gh "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
)

type SitesListViewer struct {
	Viewer               *g.GenericViewer
	LocationService      *Service
	PageSize             int
	PaginationWidgetSize int
}

func (v *SitesListViewer) List(sites []Site) gt.Node {
	var rows []g.TableRow
	for _, s := range sites {
		link := gh.A(gh.Href("/site/"+s.Id.String()),
			gh.Class(g.UrlClass),
			gt.Text(s.Name))

		rows = append(rows, g.TableRow{Values: []gt.Node{link, gt.Text(s.Group), gt.Raw(s.CreatedAt.Format("2006-01-02 / 15:04")),
			gt.Raw(s.Id.String())}})
	}

	return g.MyTable([]string{"name", "group", "created_at", "id"},
		rows)
}

func (v *SitesListViewer) makeFilteringForm(filters FilterArgs) gt.Node {

	return gh.Div(
		gh.Form(
			gh.Action("/sites"),
			gh.Method("GET"),
			gh.Table(
				gh.Tr(gh.Td(gh.Input(
					gh.Type("text"),
					gh.ID("group"),
					gh.Name("group"),
					gh.Value(filters.GetGroup()),
					gh.Placeholder("Filter by group..."),
					gh.Class("px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"),
				))),
				gh.Tr(gh.Td(gh.Input(
					gh.Type("text"),
					gh.ID("collection"),
					gh.Name("collection"),
					gh.Value(filters.GetCollection()),
					gh.Placeholder("Filter by collection..."),
					gh.Class("px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"),
				)),
					gh.Td(gh.Button(
						gh.Type("submit"),
						gh.Class("bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 px-4 rounded-md transition"),
						gt.Text("Go"),
					)))),
		),
	)

}

func (v *SitesListViewer) Render(ctx context.Context, filters FilterArgs, page int64) gt.Node {
	title := "Sites"

	filteringForm := v.makeFilteringForm(filters)

	sites, pagination, err := v.LocationService.List(ctx, filters,
		SiteAlphabeticalOrdering,
		g.PaginationParams{Page: page,
			PageSize: v.PageSize})
	if err != nil {
		return v.Viewer.OopsPage(ctx, err.Error())
	}

	baseURL, _ := url.Parse("/sites")
	values := baseURL.Query()
	values.Set("group", filters.GetGroup())
	values.Set("collection", filters.GetCollection())
	baseURL.RawQuery = values.Encode()

	paginatedViewer, err := g.MakePaginationViewer(*pagination, *baseURL, v.PageSize,
		v.PaginationWidgetSize)
	if err != nil {
		return v.Viewer.OopsPage(ctx, err.Error())
	}
	paginatedViewerPanel := gt.Group([]gt.Node{
		filteringForm,
		*paginatedViewer,
		v.List(sites),
		*paginatedViewer})

	body := gh.Div(gh.Class("px-2 md:px-4 lg:px8 py-2 md:py-4"), paginatedViewerPanel)

	return v.Viewer.BasePage(ctx, title,
		g.Head(),
		body,
	)
}
func (v *SitesListViewer) getFilterArgs(r *http.Request) FilterArgs {
	filters := FilterArgs{}
	groupNameStr := r.URL.Query().Get("group")
	if groupNameStr != "" {
		filters.Group = &groupNameStr
	}
	collectionNameStr := r.URL.Query().Get("collection")
	if collectionNameStr != "" {
		filters.Collection = &collectionNameStr
	}

	return filters
}

func (v *SitesListViewer) Handler() http.HandlerFunc {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (gt.Node, error) {

		filters := v.getFilterArgs(r)
		return v.Render(r.Context(), filters, g.GetPage(*r)), nil
	})
}
