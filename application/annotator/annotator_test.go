package annotator

import (
	"testing"

	a "github.com/lejeunel/go-image-annotator-v2/entities/annotation"
	an "github.com/lejeunel/go-image-annotator-v2/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator-v2/entities/collection"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	lbl "github.com/lejeunel/go-image-annotator-v2/entities/label"
	addbox "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/add-bbox"
	updbox "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/modify-bbox"
	del "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/remove"
)

func createAnnotator() (*Annotator, *im.Image, *FakePresenter, *FakeScroller) {
	scroller := &FakeScroller{}
	presenter := &FakePresenter{}
	image := im.NewImage(im.NewImageId(),
		*clc.NewCollection(clc.NewCollectionId(), "name"))
	label := lbl.NewLabel(lbl.NewLabelId(), "a-label")
	box := a.NewBoundingBox(a.NewAnnotationId(), 1, 1, 1, 1, *label)
	image.AddBoundingBox(*box)
	annotator := &Annotator{imageReader: &FakeImageReader{Return: image},
		scroller: scroller}
	return annotator, image, presenter, scroller

}

func TestInitializeScrollerOnStart(t *testing.T) {
	a, image, p, scroller := createAnnotator()
	a.Start(image.Id, "a-collection", p)
	if !scroller.IsInit {
		t.Fatal("expected to initialize scroller")
	}
}

func TestUpdateScrollerOnStart(t *testing.T) {
	a, image, p, _ := createAnnotator()
	a.Start(image.Id, "a-collection", p)
	if !p.UpdatedScroller {
		t.Fatal("expected to update scroller state")
	}
}
func TestPresentImageOnStart(t *testing.T) {
	a, image, p, _ := createAnnotator()
	a.Start(image.Id, "a-collection", p)
	if p.PresentedImage.Id != image.Id {
		t.Fatalf("expected to present image with id %v, got %v",
			image.Id, p.PresentedImage.Id)
	}
}

func TestAddBox(t *testing.T) {
	p := &FakePresenter{}
	adder := &FakeBoxAdder{}
	a := &Annotator{boxAdder: adder}
	imageId := im.NewImageId()
	a.AddBox(addbox.Request{ImageId: imageId}, p)
	if adder.Got.ImageId != imageId {
		t.Fatalf("expected to add bbox on image %v, got %v",
			imageId, adder.Got.ImageId)
	}
}

func TestUpdateBox(t *testing.T) {
	p := &FakePresenter{}
	updater := &FakeBoxUpdater{}
	a := &Annotator{boxUpdater: updater}
	annotationId := an.NewAnnotationId()
	a.UpdateBox(updbox.Request{AnnotationId: annotationId}, p)
	if updater.Got.AnnotationId != annotationId {
		t.Fatalf("expected to modify bbox with id %v, got %v",
			annotationId, updater.Got.AnnotationId)
	}
}

func TestRemoveBox(t *testing.T) {
	p := &FakePresenter{}
	deleter := &FakeAnnotationDeleter{}
	a := &Annotator{annotationDeleter: deleter}
	annotationId := an.NewAnnotationId()
	a.DeleteAnnotation(del.Request{Id: annotationId}, p)
	if deleter.Got.Id != annotationId {
		t.Fatalf("expected to delete bbox with id %v, got %v",
			annotationId, deleter.Got.Id)
	}
}
