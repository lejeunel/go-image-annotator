package add

type OutputPort interface {
	Error(error)
	SuccessAddMetadata()
}
