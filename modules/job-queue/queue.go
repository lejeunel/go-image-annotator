package job_queue

type JobQueue struct{}

func NewJobQueue() JobQueue {
	return JobQueue{}
}

func (q JobQueue) Submit(f func()) {
	go f()
}
