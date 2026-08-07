package list

import (
	"testing"

	ta "github.com/lejeunel/go-image-annotator/entities/task"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
	st "github.com/lejeunel/go-image-annotator/shared/testing"
	"github.com/stretchr/testify/assert"
)

func TestHandleInternalErrOnCount(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.EventLogger{ErrOnCount: e.ErrInternal})
	itr.Execute(t.Context(), pa.PaginationParams{Page: 1, PageSize: 1}, p)
	assert.Equal(t, p.GotInternalErr, true)
	assert.Equal(t, p.GotSuccess, false)
}

func TestInvalidPageShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.EventLogger{})
	itr.Execute(t.Context(), pa.PaginationParams{Page: -1}, p)
	assert.Equal(t, p.GotValidationErr, true)
	assert.Equal(t, p.GotSuccess, false)
}

func TestHandleInternalErrOnList(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.EventLogger{ErrOnList: e.ErrInternal})
	itr.Execute(t.Context(), pa.PaginationParams{Page: 1, PageSize: 1}, p)
	assert.Equal(t, p.GotInternalErr, true)
	assert.Equal(t, p.GotSuccess, false)
}

func TestListTasks(t *testing.T) {
	count := int64(2)
	pageSize := 2
	page := int64(1)

	tasks := []ta.Task{
		{Id: ta.NewTaskId(), Type: ta.CollectionCloneTask, Issuer: "user@mail.com"},
		{Id: ta.NewTaskId(), Type: ta.CollectionDeleteTask, Issuer: "another-user@mail.com"},
	}
	logger := &fk.EventLogger{Count_: count, ReturnTasks: tasks}
	p := &FakePresenter{}
	itr := New(logger)
	req := pa.PaginationParams{PageSize: pageSize, Page: page}
	itr.Execute(st.CreateCtxWithUserId(t.Context(), "user@mail.com"), req, p)
	assert.True(t, p.GotSuccess)
	assert.NoError(t, p.GotErr)
	assert.Equal(t, pageSize, len(p.Got.Tasks))
	assert.Equal(t, count, p.Got.Pagination.TotalRecords)
	assert.Equal(t, 1, int(p.Got.Pagination.TotalPages))
	assert.Equal(t, page, p.Got.Pagination.Page)
	assert.Equal(t, pageSize, p.Got.Pagination.PageSize)
}
