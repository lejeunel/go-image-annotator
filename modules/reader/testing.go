package reader

import (
	_ "embed"
)

//go:embed sample-image.jpg
var testJPGImage []byte

//go:embed sample-image.png
var testPNGImage []byte
