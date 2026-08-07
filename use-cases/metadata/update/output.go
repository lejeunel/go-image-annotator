package update

type OutputPort interface {
	Error(error)
	SuccessUpdateMetadata()
}
