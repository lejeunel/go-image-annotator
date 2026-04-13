package annotator

import (
	"testing"

	"github.com/lejeunel/go-image-annotator-v2/application/scroller"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
)

func TestScrollerButtonsWithNoPrevImage(t *testing.T) {
	v := &FakeView{}
	p := AnnotatorPresenter{v}
	p.UpdateScroller(scroller.ScrollerState{})
	if v.GotScrollerButtons.Prev.IsActive {
		t.Fatal("expected to have prev button inactive")
	}
}

func TestScrollerButtonsNextIdAndCollection(t *testing.T) {
	v := &FakeView{}
	p := AnnotatorPresenter{v}
	id := im.NewImageId()
	collection := "my-collection"
	p.UpdateScroller(scroller.ScrollerState{Next: &im.BaseImage{ImageId: id, Collection: collection}})
	if v.GotScrollerButtons.Next.ImageId != id.String() {
		t.Fatalf("expected to have next id %v, got %v", id, v.GotScrollerButtons.Next.ImageId)
	}
	if v.GotScrollerButtons.Next.Collection != collection {
		t.Fatalf("expected to have next collection %v, got %v", collection, v.GotScrollerButtons.Next.Collection)
	}
}
