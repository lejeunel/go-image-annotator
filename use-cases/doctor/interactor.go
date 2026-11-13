package doctor

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"iter"
	"sync"
	"time"

	"github.com/lejeunel/go-image-annotator/config"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	fs "github.com/lejeunel/go-image-annotator/modules/file-store"
	pv "github.com/lejeunel/go-image-annotator/modules/password-validator"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"golang.org/x/sync/errgroup"
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
	ImageFetchIdPageSize int
	ImageCheckNumWorkers int
}

func New(cfgp ConfigProvider, pvp PasswordValidatorProvider, fsp FileStoreProvider,
	dbdp DBStoreProvider, imsp ImageStoreProvider,
) Interactor {
	return Interactor{
		pvp, cfgp, fsp, dbdp, imsp, 10, 16,
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
		checkImageStoreIntegrity(
			ctx,
			imageStore,
			i.ImageFetchIdPageSize,
			i.ImageCheckNumWorkers,
			out,
		)
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

func checkImageStoreIntegrity(ctx context.Context, store ImageStore, pageSize int,
	imageCheckWorkers int,
	out OutputPort,
) {
	err := checkImageFiles(ctx, store, pageSize, imageCheckWorkers, out)
	out.DoneCheckImageIntegrity()
	if err != nil {
		out.PresentExam(Response{Description: "Image store integrity", Error: err})
	}
}

// checkImageFiles reports, for each image, whether its file exists in the
// store. Workers report under a lock, so that out is never called
// concurrently.
func checkImageFiles(ctx context.Context, store ImageStore, pageSize, workers int,
	out OutputPort,
) error {
	g, ctx := errgroup.WithContext(ctx)
	ids := make(chan im.ImageId, workers)

	g.Go(func() error {
		defer close(ids)
		for id, err := range store.IterateImageIDs(pageSize) {
			if err != nil {
				return err
			}
			select {
			case ids <- id:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return nil
	})

	var mu sync.Mutex
	for range workers {
		g.Go(func() error {
			for id := range ids {
				if err := ctx.Err(); err != nil {
					return err // another goroutine failed
				}
				exists, err := store.Exists(id)
				if err != nil {
					return err
				}
				mu.Lock()
				if *exists {
					out.SuccessCheckImageIntegrity(id)
				} else {
					out.ErrorCheckImageIntegrity(id)
				}
				mu.Unlock()
			}
			return nil
		})
	}
	return g.Wait()
}
