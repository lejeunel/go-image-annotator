package delete

type OutputPort interface {
	Error(error)
	SuccessDeleteLabel(string)
}
