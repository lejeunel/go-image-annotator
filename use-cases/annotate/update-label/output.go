package update_label

type OutputPort interface {
	Error(error)
	SuccessUpdateLabel(Response)
}
