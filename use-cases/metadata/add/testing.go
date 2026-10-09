package add

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	GotSuccess bool
	Got        *im.Image
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessAddMetadata(image im.Image) {
	p.GotSuccess = true
	p.Got = &image
}
