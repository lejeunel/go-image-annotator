package annotator

import (
	"testing"

	an "github.com/lejeunel/go-image-annotator-v2/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator-v2/entities/collection"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	lbl "github.com/lejeunel/go-image-annotator-v2/entities/label"
	add "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/add-bbox"
	upd "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/modify-bbox"
	del "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/remove"
)

func createAnnotator() (*Annotator, *im.Image, *FakeScroller, *FakeView) {
	view := &FakeView{}
	scroller := &FakeScroller{}
	image := im.NewImage(im.NewImageId(),
		*clc.NewCollection(clc.NewCollectionId(), "name"))
	label := lbl.NewLabel(lbl.NewLabelId(), "a-label")
	box := an.NewBoundingBox(an.NewAnnotationId(), 1, 1, 1, 1, *label)
	image.AddBoundingBox(*box)
	annotator := &Annotator{
		imageReader:  &FakeImageReader{Return: image},
		labelFetcher: &FakeLabelFetcher{},
		boxAdder:     &FakeBoxAdder{Returns: *box},
		boxUpdater:   &FakeBoxUpdater{Returns: &upd.Response{AnnotationId: box.Id}},
		deleter:      &FakeAnnotationDeleter{Returns: del.Response{Id: box.Id}},
		scroller:     scroller}
	return annotator, image, scroller, view

}
func TestInitializeScrollerOnStart(t *testing.T) {
	a, image, scroller, view := createAnnotator()
	a.Init(image.Id, "a-collection", view)
	if !scroller.IsInit {
		t.Fatal("expected to initialize scroller")
	}
}

func TestDrawScrollerOnStart(t *testing.T) {
	a, image, _, view := createAnnotator()
	a.Init(image.Id, "a-collection", view)
	if view.GotScrollerButtons == nil {
		t.Fatal("expected to draw scroller buttons on start")
	}

}

func TestFetchAllLabelsOnStart(t *testing.T) {
	a, image, _, view := createAnnotator()
	a.Init(image.Id, "a-collection", view)
	if view.GotLabels == nil {
		t.Fatal("expected to draw label list")
	}
}

func TestDrawImageOnStart(t *testing.T) {
	a, image, _, view := createAnnotator()
	a.Init(image.Id, "a-collection", view)
	if view.GotImage.Id != image.Id.String() {
		t.Fatalf("expected to present image with id %v, got %v",
			image.Id, view.GotImage.Id)
	}
}

func TestAddBoxShouldDraw(t *testing.T) {
	a, image, _, view := createAnnotator()
	a.AddBox(add.Request{}, view)
	box := image.BoundingBoxes[0]
	if view.GotBox.Label != box.Label.Name {
		t.Fatalf("expected to draw bbox with label %v, got %v",
			box.Label.Name, view.GotBox.Label)
	}
}

func TestRemoveBox(t *testing.T) {
	a, _, _, view := createAnnotator()
	a.DeleteAnnotation(del.Request{}, view)
	got := view.RemovedAnnotationId
	if got == nil {
		t.Fatal("expected to remove annotation")
	}
}

func TestUpdateBox(t *testing.T) {
	a, _, _, view := createAnnotator()
	a.UpdateBox(upd.Request{}, view)
	got := view.UpdatedBoxId
	if got == nil {
		t.Fatal("expected to update annotation")
	}
}
