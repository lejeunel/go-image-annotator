package presenters

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
	v "github.com/lejeunel/go-image-annotator/modules/annotator/view"
)

func MakeScrollerButtons(prev *im.BaseImage, next *im.BaseImage) v.ScrollerButtons {
	buttons := v.ScrollerButtons{}
	if next != nil {
		buttons.Next = v.ScrollerButton{
			IsActive:   true,
			Text:       "Next",
			ImageId:    next.ImageId.String(),
			Collection: next.Collection,
		}
	}
	if prev != nil {
		buttons.Prev = v.ScrollerButton{
			IsActive:   true,
			Text:       "Previous",
			ImageId:    prev.ImageId.String(),
			Collection: prev.Collection,
		}
	}
	return buttons
}
