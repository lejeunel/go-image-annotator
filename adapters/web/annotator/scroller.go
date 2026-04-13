package annotator

import (
	"fmt"
	a "github.com/lejeunel/go-image-annotator-v2/application/annotator"
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

func (p *ScrollerView) Render(buttons a.ScrollerButtons) Node {
	prevButton := myhtml.MakePreviousButton("#", false)
	nextButton := myhtml.MakeNextButton("#", false)
	if buttons.Prev.IsActive {
		prevButton = myhtml.MakePreviousButton(MakeLink(buttons.Prev.ImageId, buttons.Prev.Collection), true)
	}
	if buttons.Next.IsActive {
		prevButton = myhtml.MakeNextButton(MakeLink(buttons.Next.ImageId, buttons.Next.Collection), true)
	}
	return Table(Tr(
		Td(prevButton),
		Td(nextButton),
	))

}
