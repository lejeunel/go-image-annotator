package presenters

import (
	scr "github.com/lejeunel/go-image-annotator-v2/app/annotator/scroller"
	v "github.com/lejeunel/go-image-annotator-v2/app/annotator/view"
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
