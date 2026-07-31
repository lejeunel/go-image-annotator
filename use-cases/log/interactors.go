package log

import (
	"github.com/lejeunel/go-image-annotator/use-cases/log/find"
	"github.com/lejeunel/go-image-annotator/use-cases/log/list"
)

type Interactors struct {
	ListTasks list.Interactor
	FindTask  find.Interactor
}
