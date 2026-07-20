package fake

type JobQueue struct{}

func (q *JobQueue) Submit(f func()) {
	f()
}
