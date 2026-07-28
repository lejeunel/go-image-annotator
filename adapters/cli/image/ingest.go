package image

import (
	"context"
	"os"
	"path/filepath"

	cli "github.com/lejeunel/go-image-annotator/adapters/cli"
	s "github.com/lejeunel/go-image-annotator/app/sqlite"
	"github.com/lejeunel/go-image-annotator/config"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	ingm "github.com/lejeunel/go-image-annotator/modules/ingester"
	l "github.com/lejeunel/go-image-annotator/shared/logging"
)

type IngestPresenter struct {
	cli.ErrorPresenter
}

func (p *IngestPresenter) Success(r ingm.Response) {
	p.Info("ingested image", "id", r.ImageId)
}

func IngestDirectory(ctx context.Context, dir, collection string) {

	entries, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	app := s.NewSQLiteApp(config.Parse(), auth.NewVoidAuth(), *l.NewCliLogger())
	for _, entry := range entries {
		if !entry.IsDir() {
			f, err := os.Open(filepath.Join(dir, entry.Name()))
			if err != nil {
				return
			}
			app.Itrs.Image.Ingest.Execute(ctx, ingm.Request{Collection: collection, Reader: f},
				&IngestPresenter{cli.NewErrorPresenter()})
		}
	}
}
