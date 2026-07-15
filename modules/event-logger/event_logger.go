package event_logger

import (
	t "github.com/lejeunel/go-image-annotator/entities/task"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	"sync"
	"time"
)

type Event struct {
	Timestamp time.Time
	TaskId    t.TaskId
	Issuer    u.UserId
	Type      t.TaskType
	State     t.TaskState
	Extra     map[string]string
}

type EventLogger struct {
	mu     sync.RWMutex
	events []Event
}

func New(maxNumEvents int) *EventLogger {
	return &EventLogger{events: make([]Event, 0, maxNumEvents)}
}

func (l *EventLogger) Log(e Event) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.events = append(l.events, e)
}

// AllLIFO returns a copy of events newest-first, safe to range over
// without holding the lock.
func (l *EventLogger) AllLIFO() []Event {
	l.mu.RLock()
	n := len(l.events)
	out := make([]Event, n)
	for i, e := range l.events {
		out[n-1-i] = e // reverse while copying, single pass
	}
	l.mu.RUnlock()
	return out
}
