package builders

import (
	"embed"
)

//go:embed templates/*
var templatesFiles embed.FS

//go:embed components/*
var componentsFiles embed.FS
