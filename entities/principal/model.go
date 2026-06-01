package principal

type Principal struct {
	Id     string
	Email  string
	Groups []string
	Roles  []string
}
