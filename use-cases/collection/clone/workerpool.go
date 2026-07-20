package clone

type WorkerPool interface {
	Submit(func())
}
