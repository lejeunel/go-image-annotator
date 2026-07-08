package reset_forgotten_password

type OutputPort interface {
	Success()
	Error(error)
}
