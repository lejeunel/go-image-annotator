package assign_label

type OutputPort interface {
	SuccessAddLabel(Response)
	Error(error)
}
