package update

type OutputPort interface {
	SuccessUpdateLabel(Response)
	Error(error)
}
