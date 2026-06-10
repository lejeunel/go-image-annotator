package set_admin

type OutputPort interface {
	Success(Response)
	Error(error)
}
