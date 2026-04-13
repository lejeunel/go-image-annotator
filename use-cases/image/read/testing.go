package read

import (
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	t "github.com/lejeunel/go-image-annotator-v2/shared/testing"
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
