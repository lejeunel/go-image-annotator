package doctor

type OutputPort interface {
	SuccessMeetTheDoctor(Response)
}
