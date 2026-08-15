package scroll

type OutputPort interface {
	SuccessScroll(Response)
	Error(error)
}
