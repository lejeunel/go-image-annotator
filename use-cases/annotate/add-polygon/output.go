package add_polygon

type OutputPort interface {
	Error(error)
	SuccessAddPolygon(Response)
}
