package presenters

import (
	"testing"

	"github.com/lejeunel/go-image-annotator-v2/application/scroller"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
)

func TestScrollerButtonsWithNoPrevImage(t *testing.T) {
	buttons := MakeScrollerButtons(scroller.ScrollerState{})
	if buttons.Prev.IsActive {
		t.Fatal("expected to have prev button inactive")
	}
}

func TestScrollerButtonsNextIdAndCollection(t *testing.T) {
	id := im.NewImageId()
	collection := "my-collection"
	buttons := MakeScrollerButtons(scroller.ScrollerState{Next: &im.BaseImage{ImageId: id, Collection: collection}})
	if buttons.Next.ImageId != id.String() {
		t.Fatalf("expected to have next id %v, got %v", id, buttons.Next.ImageId)
	}
	if buttons.Next.Collection != collection {
		t.Fatalf("expected to have next collection %v, got %v", collection, buttons.Next.Collection)
	}
}
