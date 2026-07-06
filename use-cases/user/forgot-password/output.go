package forgot_password

type OutputPort interface {
	Success(Response)
	Error(error)
}
