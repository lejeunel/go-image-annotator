package slice

type OutputPort interface {
	SuccessSliceImages(Response)
	Error(error)
}
