package read

type OutputPort interface {
	Error(error)
	Success(Response)
}
