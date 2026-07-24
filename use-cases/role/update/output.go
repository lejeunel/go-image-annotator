package update

type OutputPort interface {
	SuccessUpdateRole(Response)
	Error(error)
}
