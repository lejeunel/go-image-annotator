package change_password

type OutputPort interface {
	Success()
	Error(error)
}
