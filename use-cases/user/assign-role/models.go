package assign_role

type Response struct {
	Id    string
	Roles []string
}

type Request struct {
	Id   string
	Role string
}
