package create

type Request struct {
	Name        string
	Description string
	Methods     []string
}

type Response struct {
	Name        string
	Description string
	Methods     []string
}
