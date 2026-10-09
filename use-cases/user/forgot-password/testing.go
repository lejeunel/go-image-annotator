package forgot_password

import (
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	Got        Response
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessGetForgotPasswordToken(r Response) {
	p.GotSuccess = true
	p.Got = r
}
