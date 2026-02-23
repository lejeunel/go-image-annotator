package annotation_profile

import (
	ctx "context"
	lbl "datahub/domain/labels"
	g "datahub/generic"
	"net/http"
	"net/url"

	gt "maragu.dev/gomponents"
	gh "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
)

var profileStyle = "inline-flex font-medium"

type AnnotationProfilesListViewer struct {
	Viewer               *g.GenericViewer
	Service              *Service
	PageSize             int
	PaginationWidgetSize int
}

func (v *AnnotationProfilesListViewer) concatenateLabels(labels []*lbl.Label) gt.Node {
	names := make([]gt.Node, len(labels))
	for i, label := range labels {
		names[i] = gh.Span(gt.Text(label.Name),
			gh.Class("bg-gray-100 text-gray-800 text-xs font-medium px-2.5 py-0.5 rounded"))
	}

	namesUnified := gt.Group(names)
	containedNames := gh.Div(gh.Class("flex flex-wrap gap-2"), namesUnified)
	return containedNames

}

func (v *AnnotationProfilesListViewer) List(profiles []AnnotationProfile) gt.Node {
	var rows []g.TableRow
	for _, p := range profiles {
		rows = append(rows, g.TableRow{Values: []gt.Node{
			gh.Div(gt.Text(p.Name), gh.Class(profileStyle)),
			v.concatenateLabels(p.Labels),
			gt.Raw(p.Id.String())}})
	}

	return g.MyTable([]string{"name", "labels", "id"},
		rows)
}

func (v *AnnotationProfilesListViewer) Render(ctx ctx.Context, page int64) gt.Node {
	profiles, pagination, err := v.Service.List(ctx,
		g.PaginationParams{Page: page,
			PageSize: 10})
	if err != nil {
		return v.Viewer.OopsPage(ctx, err.Error())
	}

	title := "Profiles"

	baseURL, _ := url.Parse("/profiles")
	paginatedViewer, err := g.MakePaginationViewer(*pagination, *baseURL, v.PageSize,
		v.PaginationWidgetSize)
	if err != nil {
		return v.Viewer.OopsPage(ctx, err.Error())
	}
	paginatedViewerPanel := gt.Group([]gt.Node{
		*paginatedViewer,
		v.List(profiles),
		*paginatedViewer})

	body := gh.Div(gh.Class("px-2 md:px-4 lg:px8 py-2 md:py-4"), paginatedViewerPanel)
	return v.Viewer.BasePage(ctx, title, g.Head(), body)
}

func (v *AnnotationProfilesListViewer) Handler() http.HandlerFunc {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (gt.Node, error) {

		return v.Render(r.Context(), g.GetPage(*r)), nil
	})
}
