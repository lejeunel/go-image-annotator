package ingester

import (
	"archive/zip"
	"fmt"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	ii "github.com/lejeunel/go-image-annotator/modules/image-ingester"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type ImageIngester interface {
	Ingest(r ii.Request) (*ii.Response, error)
}

type ImageStore interface {
	DeleteBatch([]im.ImageId, clc.CollectionName) error
}

type ArchiveIngester struct {
	ImageIngester
	ImageStore
}

func New(is ImageStore, ii ImageIngester) ArchiveIngester {
	return ArchiveIngester{ii, is}
}

func (i ArchiveIngester) IngestArchive(r Request) (Response, error) {
	errCtx := fmt.Errorf("ingesting zip archive")

	resp := Response{Collection: r.Collection}

	zr, err := zip.NewReader(r.ReaderAt, r.Size)
	if err != nil {
		return resp, fmt.Errorf("%w: %w", errCtx, err)
	}

	var lastErr error
	for _, file := range zr.File {
		if file.FileInfo().IsDir() {
			lastErr = fmt.Errorf("%w: found directory %v: %w", errCtx, file, e.ErrValidation)
			break
		}

		reader, err := file.Open()
		if err != nil {
			lastErr = fmt.Errorf("%w: opening file %v: %w", errCtx, file, err)
			break
		}

		r, err := i.ImageIngester.Ingest(
			ii.Request{UserId: r.UserId, Collection: r.Collection, Reader: reader},
		)
		if err != nil {
			reader.Close()
			lastErr = fmt.Errorf("%w: ingesting file %v: %w", errCtx, file.Name, err)
			break
		}

		if err := reader.Close(); err != nil {
			lastErr = fmt.Errorf("%w: closing image reader for file %v: %w", errCtx, file, err)
			break
		}
		resp.ImageIds = append(resp.ImageIds, r.ImageId)
	}

	if lastErr != nil {
		if err := i.ImageStore.DeleteBatch(resp.ImageIds, resp.Collection); err != nil {
			return resp, fmt.Errorf("%w: %w: %w", errCtx, lastErr, err)
		}
		return resp, lastErr
	}

	return resp, nil
}
