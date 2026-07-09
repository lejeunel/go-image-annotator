package list

type OutputPort interface {
	SuccessListImages(Response)
	Error(error)
}
