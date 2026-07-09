package delete

type OutputPort interface {
	Error(error)
	SuccessDeleteImage(Response)
}
