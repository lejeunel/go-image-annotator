package annotator

import (
	"testing"

	"github.com/lejeunel/go-image-annotator/app/annotator/presenters"
	an "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	st "github.com/lejeunel/go-image-annotator/shared/testing"
	addbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-bbox"
	addlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/assign-label"
	updbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/modify-bbox"
	del "github.com/lejeunel/go-image-annotator/use-cases/annotate/remove"
	updlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/update-label"
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
	st.AssertEqual(t, "initialized scroller", scroller.IsInit, true)
}
func TestDrawScrollerOnInit(t *testing.T) {
	a, image, _, view := createAnnotator()
	a.Init(image.Id.String(), "a-collection", view)
	st.AssertNonNil(t, "drawn scroller buttons", view.GotScrollerButtons)

}
func TestFetchAllLabelsOnInit(t *testing.T) {
	a, image, _, view := createAnnotator()
	a.Init(image.Id.String(), "a-collection", view)
	st.AssertNonNil(t, "drawn label list", view.GotAvailableLabels)
}
func TestDrawImageOnInit(t *testing.T) {
	a, image, _, view := createAnnotator()
	a.Init(image.Id.String(), "a-collection", view)
	st.AssertEqual(t, "image id", view.GotImage.Id, image.Id.String())
	st.AssertEqual(t, "image info id", view.GotImageInfo.Id, image.Id.String())
}
func TestAddBox(t *testing.T) {
	a, image, _, view := createAnnotator()
	a.AddBox(addbox.Request{}, view)
	st.AssertEqual(t, "added box with label", image.BoundingBoxes[0].Label.Name, view.AddedBox.Label)
}
func TestUpdateLabel(t *testing.T) {
	a, _, _, view := createAnnotator()
	a.UpdateLabelOfAnnotation(updlbl.Request{}, view)
	st.AssertNonNil(t, "update label", view.UpdatedAnnotation)
}

func TestDeleteAnnotation(t *testing.T) {
	a, _, _, view := createAnnotator()
	a.DeleteAnnotation(del.Request{}, view)
	st.AssertNonNil(t, "removed annotation", view.RemovedAnnotationId)
}
func TestUpdateBox(t *testing.T) {
	a, _, _, view := createAnnotator()
	a.UpdateBox(updbox.Request{}, view)
	st.AssertNonNil(t, "updated annotation", view.UpdatedBoxId)
}
func TestDrawImageAnnotationsOnInit(t *testing.T) {
	a, image, _, view := createAnnotator()
	a.Init(image.Id.String(), "a-collection", view)
	st.AssertNonNil(t, "got annotation ids", view.GotAnnotationIds)
	st.AssertContains(t, "bbox annotation id", *view.GotAnnotationIds, image.BoundingBoxes[0].Id.String())
	st.AssertContains(t, "label annotation id", *view.GotAnnotationIds, image.Labels[0].Id.String())
}
func TestAddLabelShouldDraw(t *testing.T) {
	a, image, _, view := createAnnotator()
	a.AddLabel(addlbl.Request{}, view)
	st.AssertEqual(t, "added box with label",
		image.Labels[0].Label.Name, view.AddedImageLabel.Label)
}
