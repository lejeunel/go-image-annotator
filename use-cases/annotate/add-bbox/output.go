package add_bbox

type OutputPort interface {
	Error(error)
	SuccessAddBox(Response)
}
