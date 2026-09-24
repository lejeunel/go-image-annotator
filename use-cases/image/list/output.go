package list

type OutputPort interface {
	SuccessPaginateImages(Response)
	Error(error)
}
