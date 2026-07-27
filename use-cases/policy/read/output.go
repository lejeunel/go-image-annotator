package read

type OutputPort interface {
	Error(error)
	SuccessReadPolicy(string)
}
