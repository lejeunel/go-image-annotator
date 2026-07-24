package delete

type OutputPort interface {
	Error(error)
	SuccessDeleteGroup(string)
}
