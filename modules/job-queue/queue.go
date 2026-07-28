package job_queue

type Interface interface {
	Submit(f func())
}

type JobQueue struct{}

func NewJobQueue() JobQueue {
	return JobQueue{}
}

func (q JobQueue) Submit(f func()) {
	go f()
}
