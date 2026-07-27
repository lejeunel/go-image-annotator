package set

type OutputPort interface {
	Error(error)
	SuccessSetPolicy(string)
}
