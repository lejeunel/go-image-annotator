package ingest

import (
	ing "github.com/lejeunel/go-image-annotator/modules/archive-ingester"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakeIngester struct {
	Got ing.Request
	Err error
}

func (i *FakeIngester) IngestArchive(r ing.Request) (ing.Response, error) {
	if i.Err != nil {
		return ing.Response{}, i.Err
	}
	i.Got = r
	return ing.Response{}, nil
}

type FakePresenter struct {
	Got        Response
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessSubmitIngestArchiveTask(r Response) {
	p.Got = r
	p.GotSuccess = true
}
