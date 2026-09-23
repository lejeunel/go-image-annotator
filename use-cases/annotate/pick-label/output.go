package pick

type OutputPort interface {
	SuccessFetchLabels([]string)
	Error(error)
}
