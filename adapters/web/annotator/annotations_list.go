package annotator

import (
	"github.com/lejeunel/go-image-annotator/modules/annotator/view"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func ShortenUUID(id string) string {
	return id[:8]
}

type LabelSelector struct {
	Labels         []string
	SelectorIsOpen bool
	Selected       *string
	AnnotationId   string
}
type AnnotationsListView struct{}

func (v *AnnotationsListView) makeRegionList(boxes []view.BoundingBox, availableLabels []string) Node {
	table := RegionTable{AvailableLabels: availableLabels}
	for _, b := range boxes {
		table.AddBox(b)
	}
	return table.Build("Regions")
}

func (v *AnnotationsListView) makeImageLabelList(imageLabels []view.ImageLabel) Node {
	table := ImageLabelTable{}
	for _, l := range imageLabels {
		table.Rows = append(table.Rows, ImageLabelRow{Label: l.Label, Id: l.Id, Author: l.Author, Time: l.Time})
	}
	return table.Build()
}

func (v *AnnotationsListView) Build(boxes []view.BoundingBox, imageLabels []view.ImageLabel, availableLabels []string) Node {
	bboxTable := v.makeRegionList(boxes, availableLabels)
	imageLabelsTable := v.makeImageLabelList(imageLabels)

	fullTable := Div(
		imageLabelsTable,
		bboxTable)
	return fullTable
}
