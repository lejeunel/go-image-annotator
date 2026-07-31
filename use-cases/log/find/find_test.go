package find

import (
	ta "github.com/lejeunel/go-image-annotator/entities/task"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleErrOnFind(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.EventLogger{ErrOnFind: e.ErrInternal})
	itr.Execute(t.Context(), ta.NewTaskId().String(), p)
	assert.Equal(t, p.GotInternalErr, true)
	assert.Equal(t, p.GotSuccess, false)
}

func TestFind(t *testing.T) {
	p := &FakePresenter{}
	task := ta.NewTask(ta.NewTaskId(), "user@mail.com", ta.CollectionCloneTask)
	itr := New(&fk.EventLogger{ReturnTask: task})
	itr.Execute(t.Context(), task.Id.String(), p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, task.Id, p.Got.Id)
}
