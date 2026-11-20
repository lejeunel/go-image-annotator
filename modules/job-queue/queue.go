package job_queue

type JobQueue interface {
	Submit(f func())
}

// AsyncJobQueue runs submitted jobs in the background, at most
// maxConcurrent at a time. The backlog of waiting jobs is unbounded.
type AsyncJobQueue struct {
	slots chan struct{}
}

func NewAsyncJobQueue(maxConcurrent int) *AsyncJobQueue {
	// A zero-capacity channel would block every job forever.
	maxConcurrent = max(maxConcurrent, 1)
	return &AsyncJobQueue{slots: make(chan struct{}, maxConcurrent)}
}

func (q *AsyncJobQueue) Submit(f func()) {
	go func() {
		q.slots <- struct{}{}        // acquire a slot, waiting if all are taken
		defer func() { <-q.slots }() // release it, even if f panics
		f()
	}()
}
