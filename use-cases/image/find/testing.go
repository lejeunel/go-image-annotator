package find

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	Got        im.Image
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessReadImage(r im.Image) {
	p.GotSuccess = true
	p.Got = r
}
