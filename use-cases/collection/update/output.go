package update

type OutputPort interface {
	SuccessUpdateCollection(Response)
	Error(error)
}
