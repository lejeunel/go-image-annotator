package reset_forgotten_password

type Request struct {
	Token          string
	FirstPassword  string
	SecondPassword string
}
