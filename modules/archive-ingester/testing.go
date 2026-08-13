package ingester

import (
	ii "github.com/lejeunel/go-image-annotator/modules/image-ingester"
)

type FakeImageIngester struct {
	Err    error
	Return ii.Response
}

func (i *FakeImageIngester) Ingest(r ii.Request) (*ii.Response, error) {
	if i.Err != nil {
		return nil, i.Err
	}
	return &i.Return, nil
}
