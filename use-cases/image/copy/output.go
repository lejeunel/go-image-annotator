package copy

type OutputPort interface {
	Error(error)
	SuccessCopyImage(Response)
}
