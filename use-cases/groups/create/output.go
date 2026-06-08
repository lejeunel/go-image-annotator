package create

type OutputPort interface {
	Success(Response)
	Error(error)
}
