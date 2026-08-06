package update

type OutputPort interface {
	Error(error)
	SuccessAddMetadata()
}
