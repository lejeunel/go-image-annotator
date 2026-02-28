package docs

import (
	"bytes"
)

type Page struct {
	ShortName string
	Title     string
	Summary   string
	Weight    int
	Content   bytes.Buffer
}
