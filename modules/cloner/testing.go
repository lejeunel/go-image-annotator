package cloner

import (
	task "github.com/lejeunel/go-image-annotator/entities/task"
)

type FakeWorkerPool struct {
	Submitted *func()
}

func (p *FakeWorkerPool) Submit(fn func()) {
	p.Submitted = &fn
}

type FakeTaskFuncBuilder struct {
	GotTask *task.CloneTask
	Err     error
	Return  func()
}

func (f *FakeTaskFuncBuilder) Build(t task.CloneTask) (func(), error) {
	f.GotTask = &t
	if f.Err != nil {
		return nil, f.Err
	}
	return nil, nil
}
