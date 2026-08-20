package update

type Request struct {
	Name           string
	NewName        string
	NewDescription string
	NewGroup       *string
}

type Response struct {
	OriginalName string
	Name         string
	Description  string
}
