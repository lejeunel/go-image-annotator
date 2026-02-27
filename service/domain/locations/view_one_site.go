package locations

import (
	"context"
	g "datahub/generic"
	gp "maragu.dev/gomponents"
	gt "maragu.dev/gomponents"
	gh "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
	"net/http"
)

type CamerasListViewer struct {
	Viewer          *g.GenericViewer
	LocationService *Service
}

func (v CamerasListViewer) List(cameras []*Camera) gp.Node {
	var rows []g.TableRow
	for _, cam := range cameras {
		linkToCamera := gh.A(gh.Href("/camera/"+cam.Id.String()),
			gh.Class(g.UrlClass),
			gt.Text(cam.Name))

		rows = append(rows, g.TableRow{Values: []gp.Node{
			linkToCamera,
			gp.Raw(cam.Id.String()),
		}})
	}

	return g.MyTable([]string{"name", "id"},
		rows)

}

func (v *CamerasListViewer) Render(ctx context.Context, siteId SiteId) gp.Node {
	site, err := v.LocationService.FindSite(ctx, siteId)
	if err != nil {
		return v.Viewer.OopsPage(ctx, err.Error())
	}

	desc := g.MakeDescriptionTable(map[string]string{
		"id":         site.Id.String(),
		"name":       site.Name,
		"group":      site.Group,
		"created_at": site.CreatedAt.Format("2006-01-02 / 15:04")},
		[]string{"id", "name", "group", "created_at"})

	cameras, err := v.LocationService.ListCamerasOfSite(ctx, site)

	if err != nil {
		return v.Viewer.OopsPage(ctx, err.Error())
	}

	camerasNode := gp.Group([]gp.Node{
		desc,
		gh.H2(gh.Class("text-3xl font-bold text-gray-900 py-2"), gp.Text("Cameras")),
		v.List(cameras)})
	body := gh.Div(gh.Class("px-2 md:px-4 lg:px8 py-2 md:py-4"), camerasNode)

	title := "Site"
	return v.Viewer.BasePage(ctx, title,
		g.Head(),
		body,
	)

}

func (v *CamerasListViewer) Handler() http.HandlerFunc {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (gp.Node, error) {
		siteId, err := NewSiteIdFromString(r.PathValue("id"))
		if err != nil {
			return v.Viewer.OopsPage(r.Context(), err.Error()), nil
		}

		return v.Render(r.Context(), *siteId), nil
	})
}
