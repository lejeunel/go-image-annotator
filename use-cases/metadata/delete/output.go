package delete

type OutputPort interface {
	Error(error)
	SuccessDeleteMetadata(string)
}
