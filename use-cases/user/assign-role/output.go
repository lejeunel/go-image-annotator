package assign_role

type OutputPort interface {
	Success(Response)
	Error(error)
}
