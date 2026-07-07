package reset_password

type Request struct {
	Token          string
	FirstPassword  string
	SecondPassword string
}
