package docs

import (
	"fmt"
	"io/fs"
)

func MakeDocs() {

	builder := NewBuilder(NewGoldMarkParser())

	err := fs.WalkDir(docsFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		f, err := docsFiles.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		fmt.Printf("Parsing %v\n", path)
		builder.AddPage(f)
		return nil
	})
	pages, err := builder.Build()
	if err != nil {
		fmt.Println(err)
	}

	for _, page := range pages {
		fmt.Println(page.Content.String())
	}
}
