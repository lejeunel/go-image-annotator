package create

type OutputPort interface {
	SuccessCreateRole(Response)
	Error(error)
}
