package delete

type OutputPort interface {
	Error(error)
	SuccessDeleteProfile(string)
}
