package annotator

import (
	"fmt"

	a "github.com/lejeunel/go-image-annotator-v2/application/annotator"
	an "github.com/lejeunel/go-image-annotator-v2/entities/annotation"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	addbox "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/add-bbox"
	updbox "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/modify-bbox"
)

type Bounds struct {
	MinX float32 `json:"minX"`
	MinY float32 `json:"minY"`
	MaxX float32 `json:"maxX"`
	MaxY float32 `json:"maxY"`
}

type BoxGeometry struct {
	XTopLeft float32 `json:"x"`
	YTopLeft float32 `json:"y"`
	W        float32 `json:"w"`
	H        float32 `json:"h"`
	Bounds   Bounds  `json:"bounds"`
}

type BoxSelector struct {
	Type     string      `json:"type"`
	Geometry BoxGeometry `json:"geometry"`
}

type BoxTarget struct {
	Selector BoxSelector `json:"selector"`
}

type AnnotoriousBody struct {
	Purpose string `json:"purpose"`
	Value   string `json:"value"`
}

type AnnotoriousBoxModel struct {
	AnnotationId string            `json:"id"`
	Properties   Properties        `json:"properties"`
	Target       BoxTarget         `json:"target"`
	Bodies       []AnnotoriousBody `json:"bodies"`
}

func (b AnnotoriousBoxModel) ExtractCoordinates() BoxCoordinates {
	return BoxCoordinates{
		Xc:     b.Target.Selector.Geometry.XTopLeft + b.Target.Selector.Geometry.W/2,
		Yc:     b.Target.Selector.Geometry.YTopLeft + b.Target.Selector.Geometry.H/2,
		Width:  b.Target.Selector.Geometry.W,
		Height: b.Target.Selector.Geometry.H,
	}
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

type BoxCoordinates struct {
	Xc     float32
	Yc     float32
	Width  float32
	Height float32
}

type Request struct {
	ImageId    im.ImageId
	Collection string
}

func ToAddBoxRequest(r AnnotoriousBoxRequest) (*addbox.Request, error) {

	imageId, err := im.NewImageIdFromString(r.ImageId)
	if err != nil {
		return nil, fmt.Errorf("submitting annotation: validating imageId: %w", err)
	}

	coords := r.Annotation.ExtractCoordinates()

	return &addbox.Request{ImageId: imageId, Collection: r.Collection,
		Label: r.Label, Xc: coords.Xc, Yc: coords.Yc, Width: coords.Width, Height: coords.Height}, nil

}

func ToUpdateBoxRequest(r AnnotoriousBoxModel) (*updbox.Request, error) {

	annotationId, err := an.NewAnnotationIdFromString(r.AnnotationId)
	if err != nil {
		return nil, fmt.Errorf("submitting annotation: validating annotationId: %w", err)
	}

	coords := r.ExtractCoordinates()

	return &updbox.Request{AnnotationId: *annotationId,
		Label: r.Bodies[0].Value, Xc: coords.Xc, Yc: coords.Yc,
		Width: coords.Width, Height: coords.Height}, nil

}

func ConvertToAnnotorious(boxes []*a.BoundingBox) []AnnotoriousBoxModel {
	result := []AnnotoriousBoxModel{}
	for _, b := range boxes {
		xtopleft := b.Xc - b.Width/2
		ytopleft := b.Yc - b.Height/2
		width := b.Width
		height := b.Height
		result = append(result,
			AnnotoriousBoxModel{
				AnnotationId: b.Id,
				Properties:   Properties{Color: b.Color},
				Bodies:       []AnnotoriousBody{{Purpose: "label", Value: b.Label}},
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
