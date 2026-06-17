package presenters

import (
	"encoding/json"
	"io"
	"net/http"

	an "github.com/lejeunel/go-image-annotator/adapters/web/annotator/annotorious"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	v "github.com/lejeunel/go-image-annotator/modules/annotator/view"
	addbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-bbox"
	addpoly "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-polygon"
	addlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/assign-label"
	updbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/modify-bbox"
	updpoly "github.com/lejeunel/go-image-annotator/use-cases/annotate/modify-polygon"
	del "github.com/lejeunel/go-image-annotator/use-cases/annotate/remove"
	updlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/update-label"
)

type JSONPresenter struct {
	Err         error
	image       *v.Image
	boxes       []v.BoundingBox
	polygons    []v.Polygon
	imageLabels []v.ImageLabel
}

func NewJSONPresenter() JSONPresenter {
	return JSONPresenter{}
}

func (p JSONPresenter) SuccessReadImage(r im.Image)             {}
func (p JSONPresenter) SuccessAddLabel(r addlbl.Response)       {}
func (p JSONPresenter) SuccessAddBox(r addbox.Response)         {}
func (p JSONPresenter) SuccessAddPolygon(r addpoly.Response)    {}
func (p JSONPresenter) SuccessUpdatePolygon(r updpoly.Response) {}
func (p JSONPresenter) SuccessUpdateBox(r updbox.Response)      {}
func (p JSONPresenter) SuccessUpdateLabel(r updlbl.Response)    {}
func (p JSONPresenter) SuccessDeleteAnnotation(r del.Response)  {}
func (p *JSONPresenter) Error(err error) {
	p.Err = err
}
func (p JSONPresenter) Write(w io.Writer) {
	// TODO write http error here
}

func (p *JSONPresenter) RenderRegionAnnotationsAsJSON(w http.ResponseWriter) {
	if p.Err != nil {
		http.Error(w, p.Err.Error(), http.StatusBadRequest)
		return
	}
	boxes := an.ConvertBoxesToAnnotorious(p.boxes)
	polygons := an.ConvertPolygonsToAnnotorious(p.polygons)
	mergedRegions := make([]any, 0, len(boxes)+len(polygons))
	for _, b := range boxes {
		mergedRegions = append(mergedRegions, b)
	}
	for _, p := range polygons {
		mergedRegions = append(mergedRegions, p)
	}

	data, err := json.Marshal(mergedRegions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
