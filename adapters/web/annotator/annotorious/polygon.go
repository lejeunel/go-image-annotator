package annotorious

import (
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	v "github.com/lejeunel/go-image-annotator/modules/annotator/view"
	addpoly "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-polygon"
)

type AnnotoriousPolygonRequest struct {
	BaseAnnotoriousRequest
	Annotation AnnotoriousPolygonModel `json:"annotation"`
}

type AnnotoriousPolygonModel struct {
	AnnotationId string            `json:"id"`
	Properties   Properties        `json:"properties"`
	Target       PolygonTarget     `json:"target"`
	Bodies       []AnnotoriousBody `json:"bodies"`
}

type PolygonGeometry struct {
	Bounds Bounds       `json:"bounds"`
	Points [][2]float32 `json:"points"`
}

type PolygonSelector struct {
	Type     string          `json:"type"`
	Geometry PolygonGeometry `json:"geometry"`
}

type PolygonTarget struct {
	Selector PolygonSelector `json:"selector"`
}

func ToAddPolygonRequest(r AnnotoriousPolygonRequest) addpoly.Request {
	return addpoly.Request{ImageId: r.ImageId, Collection: r.Collection,
		Label:  r.Label,
		Points: a.Points{Coordinates: r.Annotation.Target.Selector.Geometry.Points},
	}

}

func ConvertPolygonsToAnnotorious(polygons []v.Polygon) []AnnotoriousPolygonModel {
	result := []AnnotoriousPolygonModel{}
	for _, p := range polygons {
		result = append(result, AnnotoriousPolygonModel{
			AnnotationId: p.Id,
			Properties:   Properties{Color: p.Color},
			Bodies:       []AnnotoriousBody{{Purpose: "label", Value: p.Label}},
			Target: PolygonTarget{
				PolygonSelector{
					Type: "POLYGON",
					Geometry: PolygonGeometry{
						Points: p.Points.Coordinates,
						Bounds: Bounds{
							MinX: p.Points.MinX(),
							MinY: p.Points.MinY(),
							MaxX: p.Points.MaxX(),
							MaxY: p.Points.MaxY()}}}}})
	}
	return result
}
