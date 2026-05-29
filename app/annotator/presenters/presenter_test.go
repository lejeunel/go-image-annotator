package presenters

import (
	"testing"

	scr "github.com/lejeunel/go-image-annotator/app/annotator/scroller"
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

func TestScrollerButtonsWithNoPrevImage(t *testing.T) {
	buttons := MakeScrollerButtons(scr.ScrollerState{})
	if buttons.Prev.IsActive {
		t.Fatal("expected to have prev button inactive")
	}
}

func TestScrollerButtonsNextIdAndCollection(t *testing.T) {
	id := im.NewImageId()
	collection := "my-collection"
	buttons := MakeScrollerButtons(scr.ScrollerState{Next: &im.BaseImage{ImageId: id, Collection: collection}})
	if buttons.Next.ImageId != id.String() {
		t.Fatalf("expected to have next id %v, got %v", id, buttons.Next.ImageId)
	}
	if buttons.Next.Collection != collection {
		t.Fatalf("expected to have next collection %v, got %v", collection, buttons.Next.Collection)
	}
}
