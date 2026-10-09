package reset_forgotten_password

type OutputPort interface {
	SuccessResetForgottenPassword()
	Error(error)
}
