package read

type OutputPort interface {
	Error(error)
	SuccessReadMetadata(Response)
}
