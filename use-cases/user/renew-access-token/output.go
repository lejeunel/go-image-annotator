package renew_token

type OutputPort interface {
	Success(Response)
	Error(error)
}
