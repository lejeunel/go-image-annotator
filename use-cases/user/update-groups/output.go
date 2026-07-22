package update_group

type OutputPort interface {
	SuccessUpdateGroups(Response)
	Error(error)
}
