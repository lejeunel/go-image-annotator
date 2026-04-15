package annotator

import (
	"fmt"

	a "github.com/lejeunel/go-image-annotator-v2/application/annotator"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	addbox "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/add-bbox"
)

type Bounds struct {
	MinX float64 `json:"minX"`
	MinY float64 `json:"minY"`
	MaxX float64 `json:"maxX"`
	MaxY float64 `json:"maxY"`
}

type BoxGeometry struct {
	XTopLeft float64 `json:"x"`
	YTopLeft float64 `json:"y"`
	W        float64 `json:"w"`
	H        float64 `json:"h"`
	Bounds   Bounds  `json:"bounds"`
}

type BoxSelector struct {
	Type     string      `json:"type"`
	Geometry BoxGeometry `json:"geometry"`
}

type BoxTarget struct {
	Selector BoxSelector `json:"selector"`
}

type AnnotoriousBoxModel struct {
	AnnotationId string     `json:"id"`
	Properties   Properties `json:"properties"`
	Target       BoxTarget  `json:"target"`
}

type Properties struct {
	Color string `json:"color"`
	Label string `json:"label"`
}

type AnnotoriousBoxRequest struct {
	ImageId    string              `json:"image_id"`
	Collection string              `json:"collection"`
	Annotation AnnotoriousBoxModel `json:"annotation"`
	Label      string              `json:"label"`
}

type Request struct {
	ImageId    im.ImageId
	Collection string
}

func ConvertFromAnnotorious(r AnnotoriousBoxRequest) (*addbox.Request, error) {
	xc := r.Annotation.Target.Selector.Geometry.XTopLeft + r.Annotation.Target.Selector.Geometry.W/2
	yc := r.Annotation.Target.Selector.Geometry.YTopLeft + r.Annotation.Target.Selector.Geometry.H/2
	width := r.Annotation.Target.Selector.Geometry.W
	height := r.Annotation.Target.Selector.Geometry.H

	imageId, err := im.NewImageIdFromString(r.ImageId)
	if err != nil {
		return nil, fmt.Errorf("submitting annotation: validating imageId: %w", err)
	}

	return &addbox.Request{ImageId: imageId, Collection: r.Collection,
		Label: r.Label, Xc: float32(xc), Yc: float32(yc), Width: float32(width), Height: float32(height)}, nil

}

func ConvertToAnnotorious(boxes []*a.BoundingBox) []AnnotoriousBoxModel {
	result := []AnnotoriousBoxModel{}
	for _, b := range boxes {
		xtopleft := float64(b.Xc - b.Width/2)
		ytopleft := float64(b.Yc - b.Height/2)
		width := float64(b.Width)
		height := float64(b.Height)
		result = append(result,
			AnnotoriousBoxModel{
				AnnotationId: b.Id,
				Properties:   Properties{Color: b.Color},
				Target: BoxTarget{BoxSelector{Type: "RECTANGLE",
					Geometry: BoxGeometry{XTopLeft: xtopleft,
						YTopLeft: ytopleft,
						W:        width,
						H:        height,
						Bounds: Bounds{MinX: xtopleft,
							MinY: ytopleft,
							MaxX: xtopleft + width,
							MaxY: ytopleft + height}}}}})

	}
	return result
}
