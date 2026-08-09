package ingest

import (
	ing "github.com/lejeunel/go-image-annotator/modules/ingester"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakeIngester struct {
	Got ing.BatchRequest
}

func (i *FakeIngester) IngestArchive(r ing.BatchRequest) (ing.BatchResponse, error) {
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
