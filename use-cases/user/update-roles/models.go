package update_role

type Response struct {
	Id    string
	Roles []string
}

type Request struct {
	Id    string
	Roles []string
}
