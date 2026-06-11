package html

import (
	"embed"
)

//go:embed templates/*
var templatesFiles embed.FS
