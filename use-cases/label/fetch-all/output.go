package fetchall

type OutputPort interface {
	SuccessFetchLabels(Response)
	Error(error)
}
