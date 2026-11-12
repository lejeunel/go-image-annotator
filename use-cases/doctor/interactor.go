package doctor

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/lejeunel/go-image-annotator/config"
	fs "github.com/lejeunel/go-image-annotator/modules/file-store"
	pv "github.com/lejeunel/go-image-annotator/modules/password-validator"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

// probePayload is the content written to the file store by checkFileStore.
var probePayload = []byte("go-image-annotator doctor probe")

type FileStoreProvider interface {
	Get(config.Config) (fs.FileStore, error)
}

type ConfigProvider interface {
	Get() (*config.Config, error)
}

type PasswordValidatorProvider interface {
	Get(minEntropy int) pv.PasswordValidator
}

type Interactor struct {
	PasswordValidatorProvider
	ConfigProvider
	FileStoreProvider
}

func (i *Interactor) Execute(ctx context.Context, out OutputPort) {
	config, err := i.ConfigProvider.Get()
	parseResponse := Response{Description: "Parsing configuration"}
	if err != nil {
		parseResponse.Error = err
		out.PresentExam(parseResponse)
		return
	}
	out.PresentExam(parseResponse)

	i.examinePassword(*config, out)
	i.examineFileStore(*config, out)
}

func (i *Interactor) examinePassword(c config.Config, out OutputPort) {
	pwValidator := i.PasswordValidatorProvider.Get(c.PasswordMinEntropy)
	out.PresentExam(Response{
		Description: "Initial administrator password",
		Error:       pwValidator.Validate(c.InitialAdminPassword),
	})
}

// examineFileStore skips the sanity check when the store cannot be
// instantiated, but lets the following exams run.
func (i *Interactor) examineFileStore(c config.Config, out OutputPort) {
	store, err := i.FileStoreProvider.Get(c)
	out.PresentExam(Response{
		Description: "Instantiating file store",
		Error:       err,
	})
	if err != nil {
		return
	}

	out.PresentExam(Response{
		Description: "Writing, reading, and deleting to file store",
		Error:       checkFileStore(store),
	})
}

// checkFileStore writes probePayload under a unique key, reads it back,
// compares it and deletes it.
func checkFileStore(store fs.FileStore) error {
	key := fmt.Sprintf(".goia-doctor-%d", time.Now().UnixNano())

	if err := store.Store(key, bytes.NewReader(probePayload)); err != nil {
		return fmt.Errorf("writing: %w", err)
	}

	r, err := store.Get(key)
	if err != nil {
		return fmt.Errorf("reading: %w", err)
	}
	got, err := io.ReadAll(r)
	r.Close()
	if err != nil {
		return fmt.Errorf("reading: %w", err)
	}
	if !bytes.Equal(got, probePayload) {
		return fmt.Errorf("reading: %w", e.ErrFileStoreProbeMismatch)
	}

	if err := store.Delete(key); err != nil {
		return fmt.Errorf("deleting: %w", err)
	}
	return nil
}

type Option func(*Interactor)

func New(cfgp ConfigProvider, pvp PasswordValidatorProvider, fsp FileStoreProvider) Interactor {
	return Interactor{
		pvp, cfgp, fsp,
	}
}
