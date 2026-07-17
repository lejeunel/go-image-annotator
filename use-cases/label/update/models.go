package update

type Request struct {
	Name           string
	NewDescription string
}

type Response struct {
	Name        string
	Description string
}
