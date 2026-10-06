package update

type Request struct {
	Name           string
	NewName        string
	NewDescription *string
	NewGroup       *string
	NewProfile     *string
}

type Response struct {
	Name        string
	Description *string
	Group       *string
	Profile     *string
}
