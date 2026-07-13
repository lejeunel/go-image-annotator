package forgot_password

type Response struct {
	Email              string
	Id                 string
	PasswordResetToken string
}
