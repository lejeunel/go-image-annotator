package create

type Response struct {
	Id                  string
	PersonalAccessToken string
	IsAdmin             bool
}

type Request struct {
	Id      string
	IsAdmin bool
}
