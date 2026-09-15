package update

type OutputPort interface {
	SuccessUpdateProfile(Response)
	Error(error)
}
