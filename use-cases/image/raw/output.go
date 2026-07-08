package raw

type OutputPort interface {
	SuccessReadRawImage(Response)
	Error(error)
}
