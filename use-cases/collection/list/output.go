package list

type OutputPort interface {
	SuccessListCollections(Response)
	Error(error)
}
