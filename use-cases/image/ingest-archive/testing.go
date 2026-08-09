package ingest

import (
	ing "github.com/lejeunel/go-image-annotator/modules/image-ingester"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakeIngester struct {
	Got ing.BatchRequest
	Err error
}

func (i *FakeIngester) IngestArchive(r ing.BatchRequest) (ing.BatchResponse, error) {
	if i.Err != nil {
		return ing.BatchResponse{}, i.Err
	}
	i.Got = r
	return ing.BatchResponse{}, nil
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
