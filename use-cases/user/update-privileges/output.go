package update

type OutputPort interface {
	SuccessUpdate(Response)
	Error(error)
}
