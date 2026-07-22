package change_password

type Request struct {
	Id              string
	CurrentPassword string
	FirstPassword   string
	SecondPassword  string
}
