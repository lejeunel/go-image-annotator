package image

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	cli "github.com/lejeunel/go-image-annotator/adapters/cli"
	s "github.com/lejeunel/go-image-annotator/app/sqlite"
	"github.com/lejeunel/go-image-annotator/config"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	ingm "github.com/lejeunel/go-image-annotator/modules/ingester"
	"github.com/lejeunel/go-image-annotator/use-cases/image/ingest"
)

type IngestPresenter struct {
	cli.ErrorPresenter
}

func (p *IngestPresenter) Success(r ingm.Response) {
	fmt.Println("ingested image with id:", r.ImageId)
}

func IngestDirectory(dir, collection string) {

	entries, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	app := s.NewSQLiteApp(config.Parse(), auth.NewVoidAuth())
	for _, entry := range entries {
		ingestImage(&app.Itrs.Image.Ingest, dir, entry, collection)
	}
}

func ingestImage(itr *ingest.Interactor, dir string, entry os.DirEntry, collection string) {
	if !entry.IsDir() {
		f, err := os.Open(filepath.Join(dir, entry.Name()))
		if err != nil {
			return
		}
		itr.Execute(context.Background(), ingm.Request{Collection: collection, Reader: f}, &IngestPresenter{})
	}

}
