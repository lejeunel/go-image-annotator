package presenters

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
	v "github.com/lejeunel/go-image-annotator/modules/annotator/view"
)

func MakeScrollerButtons(adj im.AdjacentImages) v.ScrollerButtons {
	buttons := v.ScrollerButtons{}
	if adj.Next != nil {
		buttons.Next = v.ScrollerButton{
			IsActive:   true,
			Text:       "Next",
			ImageId:    adj.Next.ImageId.String(),
			Collection: adj.Next.Collection,
		}
	}
	if adj.Prev != nil {
		buttons.Prev = v.ScrollerButton{
			IsActive:   true,
			Text:       "Previous",
			ImageId:    adj.Prev.ImageId.String(),
			Collection: adj.Prev.Collection,
		}
	}
	return buttons
}
