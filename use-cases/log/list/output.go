package list

type OutputPort interface {
	SuccessListTasks(Response)
	Error(error)
}
