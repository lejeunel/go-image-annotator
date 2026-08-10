package ingest

import (
	"bytes"
	"context"
	"testing"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	ev "github.com/lejeunel/go-image-annotator/entities/event"
	grp "github.com/lejeunel/go-image-annotator/entities/group"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func Setup(t *testing.T) (Interactor, clc.Collection, grp.Group, context.Context, []byte) {
	group := grp.NewGroup(grp.NewGroupId(), "a-group")
	collection := clc.NewCollection(clc.NewCollectionId(),
		"a-collection",
		clc.WithGroup(group.Name))
	itr := New(&FakeIngester{},
		&fk.CollectionRepo{ExistingNames: []string{collection.Name}},
		&fk.FileStore{}, &fk.EventLogger{}, fk.NewLogger(), &fk.JobQueue{})
	ctx := u.AppendUserToContext(t.Context(), u.NewUser("user@mail.com"))

	data := []byte("asdf")
	return itr, collection, group, ctx, data
}

func TestHandleAuthError(t *testing.T) {
	itr, collection, _, ctx, _ := Setup(t)
	itr.CollectionRepo = &fk.CollectionRepo{
		Return: collection,
	}
	itr.Auth = &fk.Auth{Err: e.ErrAuthorization}
	p := &FakePresenter{}
	itr.Execute(ctx, Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestNonExistingCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr, _, _, ctx, _ := Setup(t)
	itr.CollectionRepo = &fk.CollectionRepo{ErrOnFind: e.ErrNotFound}
	itr.Execute(ctx, Request{}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestArchiveIsStored(t *testing.T) {
	p := &FakePresenter{}
	itr, collection, _, ctx, data := Setup(t)
	tfs := fk.FileStore{}
	itr.TemporaryFileStore = &tfs
	itr.Execute(ctx,
		Request{Reader: bytes.NewReader(data), Collection: collection.Name}, p)
	assert.Equal(t, true, bytes.Equal(tfs.GotData, data))
	assert.True(t, p.GotSuccess)
}

func TestTaskIsInitializedAndPending(t *testing.T) {
	p := &FakePresenter{}
	itr, collection, _, ctx, data := Setup(t)
	el := fk.EventLogger{}
	itr.IEventLogger = &el
	itr.Execute(ctx,
		Request{Reader: bytes.NewReader(data),
			Collection: collection.Name}, p)
	assert.True(t, el.InitializedTask)
	assert.Equal(t, ev.PendingTask, el.Events[0].State)
	assert.True(t, p.GotSuccess)
}

func TestIngestBatch(t *testing.T) {
	p := &FakePresenter{}
	itr, collection, _, ctx, data := Setup(t)
	ig := &FakeIngester{}
	itr.TemporaryFileStore = &fk.FileStore{Data: data}
	itr.ArchiveIngester = ig
	itr.Execute(ctx,
		Request{Reader: bytes.NewReader(data),
			Collection: collection.Name}, p)
	gotData := make([]byte, len(data))
	ig.Got.ReaderAt.ReadAt(gotData, 0)
	assert.True(t, bytes.Equal(gotData, data))

}

func TestArchiveIsDeleted(t *testing.T) {
	p := &FakePresenter{}
	itr, collection, _, ctx, data := Setup(t)
	tfs := &fk.FileStore{}
	itr.TemporaryFileStore = tfs
	itr.Execute(ctx,
		Request{Reader: bytes.NewReader(data),
			Collection: collection.Name}, p)
	assert.True(t, tfs.NumDeletedItems > 0)

}

func TestFailedIngestionIsLogged(t *testing.T) {
	p := &FakePresenter{}
	itr, collection, _, ctx, data := Setup(t)
	el := fk.EventLogger{}
	itr.IEventLogger = &el
	itr.ArchiveIngester = &FakeIngester{Err: e.ErrInternal}
	itr.Execute(ctx,
		Request{Reader: bytes.NewReader(data),
			Collection: collection.Name}, p)
	assert.Equal(t, ev.FailedTask, el.Events[len(el.Events)-1].State)
}

func TestIngest(t *testing.T) {
	p := &FakePresenter{}
	itr, collection, _, ctx, data := Setup(t)
	el := fk.EventLogger{}
	itr.IEventLogger = &el
	itr.ArchiveIngester = &FakeIngester{}
	itr.Execute(ctx,
		Request{Reader: bytes.NewReader(data),
			Collection: collection.Name}, p)
	assert.Equal(t, ev.DoneTask, el.Events[len(el.Events)-1].State)
}
