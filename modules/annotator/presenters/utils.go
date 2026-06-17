package presenters

import (
	"time"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	v "github.com/lejeunel/go-image-annotator/modules/annotator/view"
)

func MakeBoundingBox(b a.BoundingBox, c Colorizer) v.BoundingBox {
	res := v.BoundingBox{
		Id:     b.Id.String(),
		Label:  b.Label.Name,
		Color:  c.Colorize(b.Id.String()),
		Xc:     b.Xc,
		Yc:     b.Yc,
		Width:  b.Width,
		Height: b.Height,
		Angle:  b.Angle,
	}
	if b.Author != nil {
		res.Author = *b.Author
	} else {
		res.Author = "anonymous"
	}

	if b.Time != nil {
		t := *b.Time
		res.Time = t.Format(time.DateTime)
	}
	return res
}

func MakePolygon(p a.Polygon, c Colorizer) v.Polygon {
	res := v.Polygon{
		Id:     p.Id.String(),
		Label:  p.Label.Name,
		Color:  c.Colorize(p.Id.String()),
		Points: p.Points,
	}
	if p.Author != nil {
		res.Author = *p.Author
	} else {
		res.Author = "anonymous"
	}

	if p.Time != nil {
		t := *p.Time
		res.Time = t.Format(time.DateTime)
	}
	return res
}

func MakeImageLabels(labels []a.ImageLabel) []v.ImageLabel {
	result := []v.ImageLabel{}
	for _, l := range labels {
		row := v.ImageLabel{Id: l.Id.String(),
			Label: l.Label.Name}
		if l.Author != nil {
			row.Author = *l.Author
		} else {
			row.Author = "anonymous"
		}
		if l.Time != nil {
			t := *l.Time
			row.Time = t.Format(time.DateTime)
		}
		result = append(result, row)
	}
	return result
}

func MakeBoundingBoxes(boxes []a.BoundingBox, c Colorizer) []v.BoundingBox {
	result := []v.BoundingBox{}
	for _, b := range boxes {
		result = append(result, MakeBoundingBox(b, c))
	}
	return result
}

func MakePolygons(polygons []a.Polygon, c Colorizer) []v.Polygon {
	result := []v.Polygon{}
	for _, p := range polygons {
		result = append(result, MakePolygon(p, c))
	}
	return result
}
