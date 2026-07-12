package list

type OutputPort interface {
	SuccessListLabels(Response)
	Error(error)
}
