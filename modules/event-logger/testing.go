package event_logger

import (
	e "github.com/lejeunel/go-image-annotator/entities/event"
)

type FakeEventLogger struct {
	Got e.Event
}

func (l *FakeEventLogger) Log(event e.Event) {
	l.Got = event
}
