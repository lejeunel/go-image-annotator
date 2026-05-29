package presenters

import (
	v "github.com/lejeunel/go-image-annotator/app/annotator/view"
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
)

func MakeBoundingBox(b *a.BoundingBox, c Colorizer) *v.BoundingBox {
	return &v.BoundingBox{
		Id:     b.Id.String(),
		Label:  b.Label.Name,
		Color:  c.Colorize(b.Id.String()),
		Xc:     b.Xc,
		Yc:     b.Yc,
		Width:  b.Width,
		Height: b.Height,
	}
}

func MakeImageLabels(labels []*a.ImageLabel) []*v.ImageLabel {
	result := []*v.ImageLabel{}
	for _, l := range labels {
		result = append(result, &v.ImageLabel{Id: l.Id.String(),
			Label: l.Label.Name})
	}
	return result
}

func MakeBoundingBoxes(boxes []*a.BoundingBox, c Colorizer) []*v.BoundingBox {
	result := []*v.BoundingBox{}
	for _, b := range boxes {
		result = append(result, MakeBoundingBox(b, c))
	}
	return result
}
