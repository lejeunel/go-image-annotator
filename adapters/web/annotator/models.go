package annotator

import (
	"fmt"

	"github.com/lejeunel/go-image-annotator-v2/application/annotator/view"
	an "github.com/lejeunel/go-image-annotator-v2/entities/annotation"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	addbox "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/add-bbox"
	updbox "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/modify-bbox"
)

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

func ConvertToAnnotorious(boxes []*view.BoundingBox) []AnnotoriousBoxModel {
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
