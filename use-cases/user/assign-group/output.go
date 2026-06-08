package assign_group

type OutputPort interface {
	Success(Response)
	Error(error)
}
