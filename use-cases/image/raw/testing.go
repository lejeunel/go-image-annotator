package raw

import (
	"bytes"
	"io"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakeFileGetter struct {
	Err  error
	data []byte
}

func (f FakeFileGetter) Get(im.ImageId) (io.Reader, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	return bytes.NewBuffer(f.data), nil
}

type FakeRepo struct {
	Err         error
	ReturnSpecs *im.ImageSpecs
}

func (r FakeRepo) GetSpecs(im.ImageId) (*im.ImageSpecs, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	return r.ReturnSpecs, nil
}

type FakePresenter struct {
	Got        Response
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessReadRawImage(r Response) {
	p.GotSuccess = true
	p.Got = r
}
