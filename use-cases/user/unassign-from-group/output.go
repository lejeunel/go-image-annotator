package unassign_group

type OutputPort interface {
	Success(Response)
	Error(error)
}
