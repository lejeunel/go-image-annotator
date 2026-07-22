package list

type OutputPort interface {
	SuccessListUsers(Response)
	Error(error)
}
