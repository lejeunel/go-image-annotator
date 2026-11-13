package doctor

import (
	"errors"
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

func TestImageIntegrityDoneAfterAllImages(t *testing.T) {
	p := NewFakePresenter()
	itr := NewTestingInteractor()
	ids := []im.ImageId{im.NewImageId(), im.NewImageId()}
	itr.ImageStoreProvider = FakeImageStoreProvider{fk.ImageStore{
		IterateImageIds:      ids,
		ExistingImageFileIds: ids[:1],
	}}
	itr.Execute(t.Context(), p)
	assert.Equal(t, 1, p.NumDone)
	assert.Equal(t, 2, p.CheckedImagesAtDone)
}

func TestImageIntegrityDoneWithoutImages(t *testing.T) {
	p := NewFakePresenter()
	itr := NewTestingInteractor()
	itr.Execute(t.Context(), p)
	assert.Equal(t, 1, p.NumDone)
	assert.Equal(t, 0, p.CheckedImagesAtDone)
}

func TestImageIntegrityDoneBeforeIterationError(t *testing.T) {
	p := NewFakePresenter()
	itr := NewTestingInteractor()
	itr.ImageStoreProvider = FakeImageStoreProvider{fk.ImageStore{
		ErrOnIterateImageIDs: e.ErrInternal,
	}}
	itr.Execute(t.Context(), p)
	assert.Equal(t, 1, p.NumDone)
	assert.ErrorIs(t, p.LastExam.Error, e.ErrInternal)
	assert.Equal(t, "Image store integrity", p.LastExam.Description)
}

func TestImageIntegrityDoneBeforeExistsError(t *testing.T) {
	p := NewFakePresenter()
	itr := NewTestingInteractor()
	itr.ImageStoreProvider = FakeImageStoreProvider{fk.ImageStore{
		IterateImageIds: []im.ImageId{im.NewImageId()},
		ErrOnExists:     e.ErrInternal,
	}}
	itr.Execute(t.Context(), p)
	assert.Equal(t, 1, p.NumDone)
	assert.Equal(t, 0, p.CheckedImagesAtDone)
	assert.ErrorIs(t, p.LastExam.Error, e.ErrInternal)
}

func manyImageIds(n int) []im.ImageId {
	ids := make([]im.ImageId, n)
	for i := range ids {
		ids[i] = im.NewImageId()
	}
	return ids
}

func TestImageIntegrityChecksEveryImageConcurrently(t *testing.T) {
	p := NewFakePresenter()
	itr := NewTestingInteractor()
	ids := manyImageIds(500)
	itr.ImageStoreProvider = FakeImageStoreProvider{fk.ImageStore{
		IterateImageIds:      ids,
		ExistingImageFileIds: ids[100:],
	}}
	itr.Execute(t.Context(), p)
	assert.Equal(t, 500, p.CheckedImagesAtDone)
	assert.Equal(t, 1, p.NumDone)
	missing := 0
	for _, err := range p.Errors {
		if errors.Is(err, e.ErrNotFound) {
			missing++
		}
	}
	assert.Equal(t, 100, missing)
}

func TestImageIntegrityStopsOnExistsErrorWithManyImages(t *testing.T) {
	p := NewFakePresenter()
	itr := NewTestingInteractor()
	itr.ImageStoreProvider = FakeImageStoreProvider{fk.ImageStore{
		IterateImageIds: manyImageIds(500),
		ErrOnExists:     e.ErrInternal,
	}}
	itr.Execute(t.Context(), p)
	assert.Equal(t, 1, p.NumDone)
	assert.ErrorIs(t, p.LastExam.Error, e.ErrInternal)
}
