package update_group

type Response struct {
	Id     string
	Groups []string
}

type Request struct {
	Id     string
	Groups []string
}
