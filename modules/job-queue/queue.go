package job_queue

type JobQueue interface {
	Submit(f func())
}

type AsyncJobQueue struct{}

func NewAsyncJobQueue() AsyncJobQueue {
	return AsyncJobQueue{}
}

func (q AsyncJobQueue) Submit(f func()) {
	go f()
}
