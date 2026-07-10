package delete

type OutputPort interface {
	Error(error)
	Success(Response)
}
