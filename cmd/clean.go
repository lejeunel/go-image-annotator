package cmd

import (
	"context"
	a "datahub/app"
	c "datahub/config"
	im "datahub/domain/images"
	g "datahub/generic"
	"fmt"
	clk "github.com/jonboulle/clockwork"
	"github.com/spf13/cobra"
	"io/fs"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
)

var (
	cleanCmd = &cobra.Command{
		Use:   "clean-images",
		Short: "Clean orphan images from key-value store",
		Run: func(cmd *cobra.Command, args []string) {
			cleanOrphanImages()
		},
	}
)

func populateURIMap(uriMap *map[im.ImageId]url.URL, images []im.Image) {
	for _, image := range images {
		(*uriMap)[image.Id] = image.Uri
	}

}

func listFiles(dir string) (map[im.ImageId]string, error) {
	paths := make(map[im.ImageId]string)
	pattern := regexp.MustCompile(`^.*-?([a-f0-9]{8}-[a-f0-9]{4}-[1-5][a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12})\.(jpg|png|jpeg)$`)
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			matches := pattern.FindStringSubmatch(d.Name())
			if len(matches) > 1 {
				id, _ := im.NewImageIdFromString(matches[1])
				paths[*id] = d.Name()

			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return paths, nil
}

func cleanOrphanImages() {
	cfg := c.NewConfig()
	app, _, _, _ := a.NewApp(cfg, clk.NewRealClock(), 1)

	uris := make(map[im.ImageId]url.URL)

	ctx := context.Background()
	page := int64(1)
	pageSize := 10

	for {
		images, pagination, _ := app.Images.List(ctx, im.FilterArgs{}, im.OrderingArgs{},
			g.PaginationParams{Page: page, PageSize: pageSize}, im.FetchMetaOnly)
		populateURIMap(&uris, images)
		if pagination.Next != 0 {
			page = pagination.Next
		} else {
			break
		}
	}

	imagePath := cfg.LocalPath + "/" + "store"
	paths, err := listFiles(imagePath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("found", len(paths), "images in store")
	fmt.Println("got", len(uris), "images in database")

	for id, path := range paths {
		_, ok := uris[id]
		if !ok {
			err := os.Remove(imagePath + "/" + path)
			fmt.Println("deleted: ", path, "err:", err)
		}

	}

}
