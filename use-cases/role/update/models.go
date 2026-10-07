package update

type Request struct {
	Name           string
	NewName        string
	NewDescription *string
	NewMethods     []string
}
