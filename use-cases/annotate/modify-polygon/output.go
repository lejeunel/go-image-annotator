package modify_polygon

type OutputPort interface {
	Error(error)
	SuccessUpdatePolygon(Response)
}
