package renew_token

import (
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	Got        Response
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessRenewAPIToken(r Response) {
	p.GotSuccess = true
	p.Got = r
}
