package doctor

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"iter"
	"time"

	"github.com/lejeunel/go-image-annotator/config"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	fs "github.com/lejeunel/go-image-annotator/modules/file-store"
	pv "github.com/lejeunel/go-image-annotator/modules/password-validator"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

// probePayload is the content written to the file store by checkFileStore.
var probePayload = []byte("go-image-annotator doctor probe")

type FileStoreProvider interface {
	Get(config.Config) (fs.FileStore, error)
}

type DBStoreProvider interface {
	Get(string) (fs.FileStore, error)
}

type ConfigProvider interface {
	Get() (*config.Config, error)
}

type PasswordValidatorProvider interface {
	Get(minEntropy int) pv.PasswordValidator
}

type ImageStore interface {
	IterateImageIDs(int) iter.Seq2[im.ImageId, error]
	Exists(im.ImageId) (*bool, error)
}

type ImageStoreProvider interface {
	Get(config.Config) (ImageStore, error)
}

type Interactor struct {
	PasswordValidatorProvider
	ConfigProvider
	FileStoreProvider
	DBStoreProvider
	ImageStoreProvider
}

func New(cfgp ConfigProvider, pvp PasswordValidatorProvider, fsp FileStoreProvider,
	dbdp DBStoreProvider, imsp ImageStoreProvider,
) Interactor {
	return Interactor{
		pvp, cfgp, fsp, dbdp, imsp,
	}
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

	fileStore, err := i.FileStoreProvider.Get(*config)
	examineFileStore("file store", fileStore, err, out)

	dbDir, err := i.DBStoreProvider.Get(config.DBPath)
	examineFileStore("database store", dbDir, err, out)

	imageStore, err := i.ImageStoreProvider.Get(*config)
	if err != nil {
		out.PresentExam(Response{
			Description: "Instantiating image store",
			Error:       err,
		})
	} else {
		checkImageStoreIntegrity(imageStore, 10, out)
	}
}

func (i *Interactor) examinePassword(c config.Config, out OutputPort) {
	pwValidator := i.PasswordValidatorProvider.Get(c.PasswordMinEntropy)
	out.PresentExam(Response{
		Description: "Initial administrator password",
		Error:       pwValidator.Validate(c.InitialAdminPassword),
	})
}

// examineFileStore reports the error obtained when instantiating store, and
// probes store only when it was instantiated.
func examineFileStore(storeName string, store fs.FileStore, err error, out OutputPort) {
	out.PresentExam(Response{
		Description: fmt.Sprintf("Instantiating %s", storeName),
		Error:       err,
	})
	if err != nil {
		return
	}

	out.PresentExam(Response{
		Description: fmt.Sprintf("Writing, reading, and deleting to %s", storeName),
		Error:       checkFileStore(store),
	})
}

// checkFileStore writes a testing blob under a unique key, reads it back,
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

func checkImageStoreIntegrity(store ImageStore, pageSize int, out OutputPort) {
	response := Response{Description: "Image store integrity"}
	for id, err := range store.IterateImageIDs(pageSize) {
		if err != nil {
			response.Error = err
			out.PresentExam(response)
			return
		}
		exists, err := store.Exists(id)
		if err != nil {
			response.Error = err
			out.PresentExam(response)
			return
		}
		if !*exists {
			out.ErrorCheckImageIntegrity(id)
		} else {
			out.SuccessCheckImageIntegrity(id)
		}

	}
}
