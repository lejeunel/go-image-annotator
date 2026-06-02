package annotator

import (
	"github.com/lejeunel/go-image-annotator/app/annotator/view"
	addbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-bbox"
	updbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/modify-bbox"
)

type Request struct {
	ImageId    string
	Collection string
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
		Width: coords.Width, Height: coords.Height}, nil

}

func ConvertToAnnotorious(boxes []view.BoundingBox) []AnnotoriousBoxModel {
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
