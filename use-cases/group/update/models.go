package update

type Request struct {
	Name           string
	NewName        string
	NewDescription string
}

type Response struct {
	Name        string
	Description string
}

type Model struct {
	Name           string
	NewName        string
	NewDescription string
}
