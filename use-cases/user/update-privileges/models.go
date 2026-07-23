package update

type Response struct {
	Id     string
	Groups []string
	Roles  []string
}

type Request struct {
	Id     string
	Groups []string
	Roles  []string
}
