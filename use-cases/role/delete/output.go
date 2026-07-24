package delete

type OutputPort interface {
	Error(error)
	SuccessDeleteRole(string)
}
