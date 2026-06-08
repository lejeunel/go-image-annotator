package read

type OutputPort interface {
	Success(Response)
	Error(error)
}
