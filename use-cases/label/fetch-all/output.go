package fetchall

type OutputPort interface {
	SuccessFetchLabels([]string)
	Error(error)
}
