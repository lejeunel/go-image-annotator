package image

import (
	"context"
	"fmt"
	"github.com/lejeunel/go-image-annotator/config"
	"os"
	"path/filepath"

	cli "github.com/lejeunel/go-image-annotator/adapters/cli"
	"github.com/lejeunel/go-image-annotator/infra"
	itr "github.com/lejeunel/go-image-annotator/infra/interactors"
	ing "github.com/lejeunel/go-image-annotator/use-cases/image/ingest"
)

type IngestPresenter struct {
	cli.ErrorPresenter
}

func (p *IngestPresenter) Success(r ing.Response) {
	fmt.Println("ingested image with id:", r.ImageId)
}

func IngestDirectory(dir, collection string) {

	cfg := config.Parse()

	itr := itr.NewSQLiteImageInteractors(infra.NewSQLiteInfra(cfg.DBPath, cfg.ArtefactDir),
		cfg.AllowedImageFormats).Ingest

	entries, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		ingestImage(&itr, dir, entry, collection)
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
