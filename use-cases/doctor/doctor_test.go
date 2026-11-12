package doctor

import (
	"testing"

	"github.com/lejeunel/go-image-annotator/config"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func NewTestingInteractor() Interactor {
	return New(FakeConfigProvider{},
		&FakePasswordValidatorProvider{&fk.StringValidator{}},
		&FakeFileStoreProvider{FileStore: &fk.FileStore{}},
		&FakeDBStoreProvider{FileStore: &fk.FileStore{}},
		&FakeImageStoreProvider{&fk.ImageStore{}},
	)
}

func TestFetchingConfigFailure(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.ConfigProvider = FakeConfigProvider{Err: e.ErrConfigParsing}
	itr.Execute(t.Context(), p)
	assert.Error(t, p.LastExam.Error)
}

func TestInvalidAdminPasswordShouldFail(t *testing.T) {
	p := &FakePresenter{}
	pwval := &fk.StringValidator{Invalid: true}
	itr := NewTestingInteractor()
	itr.PasswordValidatorProvider = &FakePasswordValidatorProvider{pwval}
	itr.Execute(t.Context(), p)
	assert.True(t, p.GotErr(e.ErrValidation))
}

func TestInvalidImageStoreConfig(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.FileStoreProvider = &FakeFileStoreProvider{Err: e.ErrValidation}
	itr.Execute(t.Context(), p)
	assert.True(t, p.GotErr(e.ErrValidation))
}

func TestFileStoreRoundTripSucceeds(t *testing.T) {
	p := &FakePresenter{}
	store := &fk.FileStore{}
	itr := NewTestingInteractor()
	itr.FileStoreProvider = &FakeFileStoreProvider{FileStore: store}
	itr.Execute(t.Context(), p)
	assert.Empty(t, p.Errors)
	assert.Equal(t, probePayload, store.GotData)
	assert.Equal(t, 1, store.NumGets)
	assert.Equal(t, 1, store.NumDeletedItems)
}

func TestFileStoreWriteFailure(t *testing.T) {
	p := &FakePresenter{}
	store := &fk.FileStore{ErrOnStore: e.ErrInternal}
	itr := NewTestingInteractor()
	itr.FileStoreProvider = &FakeFileStoreProvider{FileStore: store}
	itr.Execute(t.Context(), p)
	assert.True(t, p.GotErr(e.ErrInternal))
	assert.Equal(t, 0, store.NumGets)
}

func TestFileStoreReadFailure(t *testing.T) {
	p := &FakePresenter{}
	store := &fk.FileStore{ErrOnGet: e.ErrNotFound}
	itr := NewTestingInteractor()
	itr.FileStoreProvider = &FakeFileStoreProvider{FileStore: store}
	itr.Execute(t.Context(), p)
	assert.True(t, p.GotErr(e.ErrNotFound))
}

func TestFileStoreContentMismatch(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.FileStoreProvider = &FakeFileStoreProvider{FileStore: CorruptedFileStore{&fk.FileStore{}}}
	itr.Execute(t.Context(), p)
	assert.True(t, p.GotErr(e.ErrFileStoreProbeMismatch))
}

func TestFileStoreDeleteFailure(t *testing.T) {
	p := &FakePresenter{}
	store := &fk.FileStore{ErrOnDelete: e.ErrInternal}
	itr := NewTestingInteractor()
	itr.FileStoreProvider = &FakeFileStoreProvider{FileStore: store}
	itr.Execute(t.Context(), p)
	assert.True(t, p.GotErr(e.ErrInternal))
}

func TestInvalidDBStore(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	itr.ConfigProvider = FakeConfigProvider{Config: config.Config{DBPath: "/some/db/dir"}}
	provider := &FakeDBStoreProvider{Err: e.ErrValidation}
	itr.DBStoreProvider = provider
	itr.Execute(t.Context(), p)
	assert.Equal(t, "/some/db/dir", provider.GotPath)
	assert.True(t, p.GotErr(e.ErrValidation))
}

func TestDBStoreWriteFailure(t *testing.T) {
	p := &FakePresenter{}
	store := &fk.FileStore{ErrOnStore: e.ErrInternal}
	itr := NewTestingInteractor()
	itr.DBStoreProvider = &FakeDBStoreProvider{FileStore: store}
	itr.Execute(t.Context(), p)
	assert.True(t, p.GotErr(e.ErrInternal))
}

func TestMissingImageFileShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	id := im.NewImageId()
	itr.ImageStoreProvider = FakeImageStoreProvider{
		fk.ImageStore{IterateImageIds: []im.ImageId{id}},
	}
	itr.Execute(t.Context(), p)
	assert.True(t, p.GotErr(e.ErrNotFound))
	assert.Equal(t, 1, p.CheckedImages)
}

func TestExistingImageFileShouldSucceed(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor()
	id := im.NewImageId()
	itr.ImageStoreProvider = FakeImageStoreProvider{fk.ImageStore{
		IterateImageIds:      []im.ImageId{id},
		ExistingImageFileIds: []im.ImageId{id},
	}}
	itr.Execute(t.Context(), p)
	assert.False(t, p.GotErr(e.ErrNotFound))
	assert.Equal(t, 1, p.CheckedImages)
}
