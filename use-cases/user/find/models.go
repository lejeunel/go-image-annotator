package find

type Request struct {
	Id string
}

type Response struct {
	Id     string
	Groups []string
	Roles  []string
}
