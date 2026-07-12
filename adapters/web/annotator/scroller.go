package annotator

import (
	"fmt"

	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	"github.com/lejeunel/go-image-annotator/modules/annotator/view"
	scr "github.com/lejeunel/go-image-annotator/modules/scroller"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type ScrollerView struct {
}

func MakeLink(imageId, collection string) string {
	return fmt.Sprintf("image?id=%v&collection=%v",
		imageId, collection)

}

func (p *ScrollerView) Render(buttons view.ScrollerButtons) Node {
	prevURL, nextURL := "#", "#"
	if buttons.Prev.IsActive {
		prevURL = MakeLink(buttons.Prev.ImageId, buttons.Prev.Collection)
	}
	if buttons.Next.IsActive {
		nextURL = MakeLink(buttons.Next.ImageId, buttons.Next.Collection)
	}
	return Div(Class("flex gap-2"),
		cmp.MakeNavigationButton(prevURL, buttons.Prev.IsActive, scr.ScrollPrevious, "Previous"),
		cmp.MakeNavigationButton(nextURL, buttons.Next.IsActive, scr.ScrollNext, "Next"),
	)

}
