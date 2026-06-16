package annotorious

import (
	v "github.com/lejeunel/go-image-annotator/modules/annotator/view"
	addbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-bbox"
	updbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/modify-bbox"
)

type BoxGeometry struct {
	XTopLeft float32 `json:"x"`
	YTopLeft float32 `json:"y"`
	W        float32 `json:"w"`
	H        float32 `json:"h"`
	Bounds   Bounds  `json:"bounds"`
	Rot      float32 `json:"rot"`
}

type BoxSelector struct {
	Type     string      `json:"type"`
	Geometry BoxGeometry `json:"geometry"`
}

type BoxTarget struct {
	Selector BoxSelector `json:"selector"`
}

type AnnotoriousBoxModel struct {
	AnnotationId string            `json:"id"`
	Properties   Properties        `json:"properties"`
	Target       BoxTarget         `json:"target"`
	Bodies       []AnnotoriousBody `json:"bodies"`
}

type AnnotoriousBoxRequest struct {
	BaseAnnotoriousRequest
	Annotation AnnotoriousBoxModel `json:"annotation"`
}

type BoxCoordinates struct {
	Xc     float32
	Yc     float32
	Width  float32
	Height float32
	Angle  float32
}

func (b AnnotoriousBoxModel) ExtractCoordinates() BoxCoordinates {
	return BoxCoordinates{
		Xc:     b.Target.Selector.Geometry.XTopLeft + b.Target.Selector.Geometry.W/2,
		Yc:     b.Target.Selector.Geometry.YTopLeft + b.Target.Selector.Geometry.H/2,
		Width:  b.Target.Selector.Geometry.W,
		Height: b.Target.Selector.Geometry.H,
		Angle:  b.Target.Selector.Geometry.Rot,
	}
}

func ToAddBoxRequest(r AnnotoriousBoxRequest) (*addbox.Request, error) {

	coords := r.Annotation.ExtractCoordinates()

	return &addbox.Request{ImageId: r.ImageId, Collection: r.Collection,
		Label: r.Label, Xc: coords.Xc, Yc: coords.Yc, Width: coords.Width, Height: coords.Height}, nil

}

func ToUpdateBoxRequest(r AnnotoriousBoxModel) (*updbox.Request, error) {

	coords := r.ExtractCoordinates()

	return &updbox.Request{AnnotationId: r.AnnotationId,
		Label: r.Bodies[0].Value, Xc: coords.Xc, Yc: coords.Yc,
		Width: coords.Width, Height: coords.Height, Angle: coords.Angle}, nil

}

func ConvertBoxesToAnnotorious(boxes []v.BoundingBox) []AnnotoriousBoxModel {
	result := []AnnotoriousBoxModel{}
	for _, b := range boxes {
		xtopleft := b.Xc - b.Width/2
		ytopleft := b.Yc - b.Height/2
		result = append(result,
			AnnotoriousBoxModel{
				AnnotationId: b.Id,
				Properties:   Properties{Color: b.Color},
				Bodies:       []AnnotoriousBody{{Purpose: "label", Value: b.Label}},
				Target: BoxTarget{BoxSelector{Type: "RECTANGLE",
					Geometry: BoxGeometry{
						XTopLeft: xtopleft,
						YTopLeft: ytopleft,
						W:        b.Width,
						H:        b.Height,
						Rot:      b.Angle,
						Bounds: Bounds{MinX: xtopleft,
							MinY: ytopleft,
							MaxX: xtopleft + b.Width,
							MaxY: ytopleft + b.Height}}}}})

	}
	return result
}
