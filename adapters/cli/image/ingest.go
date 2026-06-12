package image

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	cli "github.com/lejeunel/go-image-annotator/adapters/cli"
	a "github.com/lejeunel/go-image-annotator/adapters/sqlite/app"
	ing "github.com/lejeunel/go-image-annotator/use-cases/image/ingest"
)

type IngestPresenter struct {
	cli.ErrorPresenter
}

func (p *IngestPresenter) Success(r ing.Response) {
	fmt.Println("ingested image with id:", r.ImageId)
}

func IngestDirectory(dir, collection string) {

	entries, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	app := a.NewSQLiteApp()
	for _, entry := range entries {
		ingestImage(&app.Itrs.Image.Ingest, dir, entry, collection)
	}
}

func ingestImage(itr *ing.Interactor, dir string, entry os.DirEntry, collection string) {
	if !entry.IsDir() {
		f, err := os.Open(filepath.Join(dir, entry.Name()))
		if err != nil {
			fmt.Println(err)
			return
		}
		itr.Execute(context.Background(), ing.Request{Collection: collection, Reader: f}, &IngestPresenter{})
	}

}
