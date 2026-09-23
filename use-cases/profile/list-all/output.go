package list

type OutputPort interface {
	SuccessListAllProfiles([]string)
	Error(error)
}
