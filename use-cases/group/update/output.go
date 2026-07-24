package update

type OutputPort interface {
	SuccessUpdateGroup(Response)
	Error(error)
}
