package presenters

import (
	"encoding/json"
	"net/http"

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

type AnnotoriousPresenter struct {
	Colorizer
	Err      error
	image    *v.Image
	boxes    []v.BoundingBox
	polygons []v.Polygon
}

func NewAnnotoriousPresenter(colorizer Colorizer) AnnotoriousPresenter {
	return AnnotoriousPresenter{Colorizer: colorizer}
}

func (p *AnnotoriousPresenter) SuccessReadImage(r im.Image) {
	p.boxes = MakeBoundingBoxes(r.BoundingBoxes, p.Colorizer)
	p.polygons = MakePolygons(r.Polygons, p.Colorizer)
}
func (p AnnotoriousPresenter) SuccessAddLabel(r addlbl.Response)       {}
func (p AnnotoriousPresenter) SuccessAddBox(r addbox.Response)         {}
func (p AnnotoriousPresenter) SuccessAddPolygon(r addpoly.Response)    {}
func (p AnnotoriousPresenter) SuccessUpdatePolygon(r updpoly.Response) {}
func (p AnnotoriousPresenter) SuccessUpdateBox(r updbox.Response)      {}
func (p AnnotoriousPresenter) SuccessUpdateLabel(r updlbl.Response)    {}
func (p AnnotoriousPresenter) SuccessDeleteAnnotation(r del.Response)  {}
func (p *AnnotoriousPresenter) Error(err error) {
	p.Err = err
}
func (p AnnotoriousPresenter) Write(w http.ResponseWriter) {
	// TODO write http error here
}

func (p *AnnotoriousPresenter) RenderRegionAnnotationsAsJSON(w http.ResponseWriter) {
	if p.Err != nil {
		http.Error(w, p.Err.Error(), http.StatusBadRequest)
		return
	}
	boxes := ConvertBoxesToAnnotorious(p.boxes)
	polygons := ConvertPolygonsToAnnotorious(p.polygons)
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
