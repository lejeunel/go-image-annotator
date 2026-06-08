package renew_token

type Response struct {
	Id                  string
	PersonalAccessToken string
}

type Request struct {
	Id string
}
