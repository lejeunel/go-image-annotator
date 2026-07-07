package find

type OutputPort interface {
	Success(Response)
	Error(error)
}
