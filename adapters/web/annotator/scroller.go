package annotator

import (
	"fmt"

	"github.com/lejeunel/go-image-annotator-v2/application/annotator/view"
	"github.com/lejeunel/go-image-annotator-v2/application/scroller"
	myhtml "github.com/lejeunel/go-image-annotator-v2/shared/html"

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
	return Table(Tr(
		Td(myhtml.MakeNavigationButton(prevURL, buttons.Prev.IsActive, scroller.ScrollPrevious, "Previous")),
		Td(myhtml.MakeNavigationButton(nextURL, buttons.Next.IsActive, scroller.ScrollNext, "Next")),
	))

}
