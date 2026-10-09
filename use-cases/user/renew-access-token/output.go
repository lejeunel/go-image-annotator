package renew_token

type OutputPort interface {
	SuccessRenewAPIToken(Response)
	Error(error)
}
