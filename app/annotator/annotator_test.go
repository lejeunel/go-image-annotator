package annotator

import (
	"testing"

	"github.com/lejeunel/go-image-annotator/app/annotator/presenters"
	an "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	addbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-bbox"
	addlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/assign-label"
	updbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/modify-bbox"
	del "github.com/lejeunel/go-image-annotator/use-cases/annotate/remove"
	updlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/update-label"
	"github.com/stretchr/testify/assert"
)

func createAnnotator() (*Annotator, *im.Image, *FakeScroller, *FakeView) {
	view := &FakeView{}
	scroller := &FakeScroller{}
	image := im.NewImage(im.NewImageId(),
		clc.NewCollection(clc.NewCollectionId(), "name"))
	label := lbl.NewLabel(lbl.NewLabelId(), "a-label")
	box := an.NewBoundingBox(an.NewAnnotationId(), 1, 1, 1, 1, label)
	image.AddBoundingBox(*box)
	image.AddLabel(label)
	annotator := NewAnnotator(
		scroller,
		&FakeImageReader{Return: &image},
		&FakeBoxAdder{Returns: *box},
		&FakeBoxUpdater{Returns: &updbox.Response{AnnotationId: box.Id}},
		&FakeAnnotationDeleter{Returns: del.Response{Id: box.Id}},
		&FakeLabelFetcher{},
		&FakeLabelUpdater{},
		&FakeLabelAdder{Returns: addlbl.Response{ImageId: image.Id, Collection: image.Collection.Name, Label: label.Name}},
		presenters.NewPresenter())
	return annotator, &image, scroller, view

}
func TestInitializeScrollerOnStart(t *testing.T) {
	a, image, scroller, view := createAnnotator()
	a.Init(image.Id.String(), "a-collection", view)
	assert.Equal(t, scroller.IsInit, true, "initialized scroller")
}
func TestDrawScrollerOnInit(t *testing.T) {
	a, image, _, view := createAnnotator()
	a.Init(image.Id.String(), "a-collection", view)
	assert.NotNil(t, view.GotScrollerButtons, "drawn scroller buttons")

}
func TestFetchAllLabelsOnInit(t *testing.T) {
	a, image, _, view := createAnnotator()
	a.Init(image.Id.String(), "a-collection", view)
	assert.NotNil(t, view.GotAvailableLabels, "drawn label list")
}
func TestDrawImageOnInit(t *testing.T) {
	a, image, _, view := createAnnotator()
	a.Init(image.Id.String(), "a-collection", view)
	assert.Equal(t, view.GotImage.Id, image.Id.String(), "image id")
	assert.Equal(t, view.GotImageInfo.Id, image.Id.String(), "image info id")
}
func TestAddBox(t *testing.T) {
	a, image, _, view := createAnnotator()
	a.AddBox(addbox.Request{}, view)
	assert.Equal(t, image.BoundingBoxes[0].Label.Name, view.AddedBox.Label, "added box with label")
}
func TestUpdateLabel(t *testing.T) {
	a, _, _, view := createAnnotator()
	a.UpdateLabel(updlbl.Request{}, view)
	assert.NotNil(t, view.UpdatedAnnotation, "update label")
}
func TestDeleteAnnotation(t *testing.T) {
	a, _, _, view := createAnnotator()
	a.DeleteAnnotation(del.Request{}, view)
	assert.NotNil(t, view.RemovedAnnotationId, "removed annotation")
}
func TestUpdateBox(t *testing.T) {
	a, _, _, view := createAnnotator()
	a.UpdateBox(updbox.Request{}, view)
	assert.NotNil(t, view.UpdatedBoxId, "updated annotation")
}
func TestDrawImageAnnotationsOnInit(t *testing.T) {
	a, image, _, view := createAnnotator()
	a.Init(image.Id.String(), "a-collection", view)
	assert.NotNil(t, view.GotAnnotationIds, "got annotation ids")
	assert.Contains(t, *view.GotAnnotationIds, image.BoundingBoxes[0].Id.String(), "bbox annotation id")
	assert.Contains(t, *view.GotAnnotationIds, image.Labels[0].Id.String(), "label annotation id")
}
func TestAddLabelShouldDraw(t *testing.T) {
	a, image, _, view := createAnnotator()
	a.AddLabel(addlbl.Request{}, view)
	assert.Equal(t, image.Labels[0].Label.Name, view.AddedImageLabel.Label, "added box with label")
}
