package sqlite

import (
	el "github.com/lejeunel/go-image-annotator/modules/event-logger"
	l "github.com/lejeunel/go-image-annotator/use-cases/log"
	ll "github.com/lejeunel/go-image-annotator/use-cases/log/list"
)

func NewSQLiteLogInteractors(el el.Interface) l.Interactors {
	return l.Interactors{
		ListTasks: ll.New(el),
	}
}
