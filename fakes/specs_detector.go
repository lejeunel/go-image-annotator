package fake

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
	"io"
)

type SpecsDetector struct {
	Err    error
	Return im.Specs
}

func (d *SpecsDetector) Detect(r io.Reader) (*im.Specs, io.Reader, error) {
	if d.Err != nil {
		return nil, nil, d.Err
	}
	return &d.Return, r, nil
}
