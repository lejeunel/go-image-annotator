package create

type Response struct {
	Id      string
	IsAdmin bool
}

type Request struct {
	Id       string
	IsAdmin  bool
	Password *string
}
