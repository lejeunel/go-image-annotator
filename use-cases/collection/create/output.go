package create

type OutputPort interface {
	SuccessCreateCollection(Response)
	Error(error)
}
