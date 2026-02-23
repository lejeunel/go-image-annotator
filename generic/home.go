package generic

import (
	gp "maragu.dev/gomponents"
	gh "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
	"net/http"
)

type HomeViewer struct {
	Viewer *GenericViewer
}

func (v *HomeViewer) Handler() http.HandlerFunc {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (gp.Node, error) {
		return v.Viewer.BasePage(r.Context(),
			"Home",
			Head(),
			gh.P(gp.Text("Intro goes here."))), nil
	})
}
