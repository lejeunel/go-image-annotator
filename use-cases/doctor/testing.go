package doctor

type FakePresenter struct {
	Got Response
}

func (p *FakePresenter) SuccessMeetTheDoctor(r Response) {
	p.Got = r
}
