package bootstrap

type OutputPort interface {
	SuccessBootstrap(Response)
	Error(error)
}
