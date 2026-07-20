package fake

type WorkerPool struct{}

func (w *WorkerPool) Submit(f func()) {
	f()
}
