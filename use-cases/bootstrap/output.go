package bootstrap

type OutputPort interface {
	SuccessBootstrap()
	Error(error)
}
