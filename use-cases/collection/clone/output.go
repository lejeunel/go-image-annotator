package clone

type OutputPort interface {
	SuccessSubmitCloneTask(Response)
	Error(error)
}
