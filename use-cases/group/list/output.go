package list

type OutputPort interface {
	Success(Response)
	Error(error)
}
