package event_logger

import (
	e "github.com/lejeunel/go-image-annotator/entities/event"
	"sync"
)

type Interface interface {
	Log(e.Event)
}

type EventLogger struct {
	mu     sync.RWMutex
	events []e.Event
}

func New(maxNumEvents int) *EventLogger {
	return &EventLogger{events: make([]e.Event, 0, maxNumEvents)}
}

func (l *EventLogger) Log(e e.Event) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.events = append(l.events, e)
}

// AllLIFO returns a copy of events newest-first, safe to range over
// without holding the lock.
func (l *EventLogger) AllLIFO() []e.Event {
	l.mu.RLock()
	n := len(l.events)
	out := make([]e.Event, n)
	for i, e := range l.events {
		out[n-1-i] = e // reverse while copying, single pass
	}
	l.mu.RUnlock()
	return out
}
