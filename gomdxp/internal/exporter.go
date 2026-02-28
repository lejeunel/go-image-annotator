package docs

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func ParseMarkDownToHTML(builder *Builder, inputFS fs.FS) error {

	err := fs.WalkDir(inputFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		f, err := inputFS.Open(path)
		data, err := io.ReadAll(f)
		if err != nil {
			f.Close()
			return err
		}
		f.Close()

		fmt.Printf("Found %v\n", path)
		builder.AddPage(bytes.NewReader(data))
		return nil
	})
	if err != nil {
		return err
	}
	if err := builder.Parse(); err != nil {
		return err
	}
	return nil
}

func ExportPagesToHTML(builder *Builder, outputPath string) error {
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return fmt.Errorf("creating output HTML directory: %w", err)
	}
	for _, page := range builder.Build() {
		fullOutputPath := filepath.Join(outputPath, fmt.Sprintf("%v.html", page.ShortName))
		if err := os.WriteFile(fullOutputPath, page.Content.Bytes(), 0644); err != nil {
			return err
		}
	}
	return nil
}
