package annotator

import (
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	"github.com/lejeunel/go-image-annotator/modules/annotator/view"
	rt "github.com/lejeunel/go-image-annotator/routes"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type ScrollerView struct{}

func MakeLink(imageId, collection string, f im.FilterStr, o im.OrderStr) string {
	url := rt.AddQueryParams("#",
		rt.CollectionArgName, collection,
		rt.ImageIdArgName, imageId,
		rt.FilterQueryArgName, f,
		rt.OrderingQueryArgName, o)
	return url.String()
}

func (p *ScrollerView) Render(buttons view.ScrollerButtons, f im.FilterStr, o im.OrderStr) Node {
	prevURL, nextURL := "#", "#"
	if buttons.Prev.IsActive {
		prevURL = MakeLink(buttons.Prev.ImageId, buttons.Prev.Collection, f, o)
	}
	if buttons.Next.IsActive {
		nextURL = MakeLink(buttons.Next.ImageId, buttons.Next.Collection, f, o)
	}
	return Div(Class("flex gap-2 mr-4"),
		cmp.MakeNavigationButton(prevURL, buttons.Prev.IsActive, im.ScrollPrevious, "Previous"),
		cmp.MakeNavigationButton(nextURL, buttons.Next.IsActive, im.ScrollNext, "Next"),
	)
}
