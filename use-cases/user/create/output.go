package create

type OutputPort interface {
	SuccessCreateUser(Response)
	Error(error)
}
