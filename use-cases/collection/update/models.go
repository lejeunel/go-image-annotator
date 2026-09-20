package update

type Request struct {
	Name           string
	NewName        string
	NewDescription string
	NewGroup       *string
	NewProfile     *string
}

type Response struct {
	OriginalName string
	Name         string
	Description  string
}
