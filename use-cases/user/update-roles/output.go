package update_role

type OutputPort interface {
	SuccessUpdateRoles(Response)
	Error(error)
}
