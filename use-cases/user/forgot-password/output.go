package forgot_password

type OutputPort interface {
	SuccessGetForgotPasswordToken(Response)
	Error(error)
}
