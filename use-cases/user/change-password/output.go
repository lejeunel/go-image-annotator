package change_password

type OutputPort interface {
	SuccessChangePassword()
	Error(error)
}
