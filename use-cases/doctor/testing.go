package doctor

type FakePresenter struct {
	Got Response
}

func (p *FakePresenter) PresentExam(r Response) {
	p.Got = r
}
