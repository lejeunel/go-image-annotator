package list

type OutputPort interface {
	SuccessListProfiles(Response)
	Error(error)
}
