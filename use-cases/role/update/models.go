package update

type Request struct {
	Name           string
	NewName        string
	NewDescription string
	NewMethods     []string
}

type Response struct {
	Name        string
	Description string
	Methods     []string
}
