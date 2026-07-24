package create

type Response struct {
	Id     string
	Groups []string
	Roles  []string
}

type Request struct {
	Id       string
	Roles    []string
	Groups   []string
	Password *string
}
