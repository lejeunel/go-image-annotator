package delete

type OutputPort interface {
	Error(error)
	SuccessDeleteCollection(Response)
}
