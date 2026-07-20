package event

import (
	"fmt"
	"time"
)

type State string

const (
	PendingTask State = "pending"
	StartedTask State = "started"
	FailedTask  State = "failed"
	DoneTask    State = "done"
)

func (r State) String() string {
	return string(r)
}

func (r State) Valid() bool {
	switch r {
	case PendingTask, StartedTask, FailedTask, DoneTask:
		return true
	default:
		return false
	}
}

func ParseState(s string) (State, error) {
	r := State(s)
	if !r.Valid() {
		return "", fmt.Errorf("invalid state %q", s)
	}
	return r, nil
}

type Event struct {
	Time  time.Time
	State State
	Extra map[string]string
	Error string
}

func (e *Event) SetState(now time.Time, state State) Event {
	e.State = state
	e.Time = now
	return *e
}

func (e *Event) SetError(now time.Time, err string) Event {
	e.Error = err
	e.Time = now
	return *e
}

func (e *Event) SetExtra(now time.Time, extra map[string]string) Event {
	e.Extra = extra
	e.Time = now
	return *e
}

func New() Event {
	return Event{
		Extra: make(map[string]string)}

}
