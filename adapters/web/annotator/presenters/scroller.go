package presenters

import (
	v "github.com/lejeunel/go-image-annotator/modules/annotator/view"
	scr "github.com/lejeunel/go-image-annotator/modules/scroller"
)

func MakeScrollerButtons(s scr.ScrollerState) v.ScrollerButtons {
	buttons := v.ScrollerButtons{}
	if s.Next != nil {
		buttons.Next = v.ScrollerButton{IsActive: true,
			Text:       "Next",
			ImageId:    s.Next.ImageId.String(),
			Collection: s.Next.Collection}
	}
	if s.Previous != nil {
		buttons.Prev = v.ScrollerButton{IsActive: true,
			Text:       "Previous",
			ImageId:    s.Previous.ImageId.String(),
			Collection: s.Previous.Collection}
	}
	return buttons
}
