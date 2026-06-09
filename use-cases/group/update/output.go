package update

type OutputPort interface {
	Success(Response)
	Error(error)
}
