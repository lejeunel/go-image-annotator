package presenters

import (
	"testing"

	im "github.com/lejeunel/go-image-annotator/entities/image"
)

func TestScrollerButtonsWithNoPrevImage(t *testing.T) {
	buttons := MakeScrollerButtons(im.AdjacentImages{Next: &im.BaseImage{}})
	if buttons.Prev.IsActive {
		t.Fatal("expected to have prev button inactive")
	}
}

func TestScrollerButtonsNextIdAndCollection(t *testing.T) {
	id := im.NewImageId()
	collection := "my-collection"
	buttons := MakeScrollerButtons(
		im.AdjacentImages{
			Next: &im.BaseImage{ImageId: id, Collection: collection},
		},
	)
	if buttons.Next.ImageId != id.String() {
		t.Fatalf("expected to have next id %v, got %v", id, buttons.Next.ImageId)
	}
}
