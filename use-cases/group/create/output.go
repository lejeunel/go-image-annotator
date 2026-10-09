package create

type OutputPort interface {
	SuccessCreateGroup(Response)
	Error(error)
}
