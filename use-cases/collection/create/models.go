package create

type Response struct {
	Name        string
	Group       string
	Description string
}

type Request struct {
	Name        string
	Description string
	Group       *string
}
