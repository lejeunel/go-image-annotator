package web

import (
	"embed"
)

//go:embed templates/*
var templatesFiles embed.FS
