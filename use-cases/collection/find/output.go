package find

type OutputPort interface {
	Error(error)
	Success(Response)
}
