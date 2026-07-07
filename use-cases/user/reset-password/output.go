package reset_password

type OutputPort interface {
	Success()
	Error(error)
}
