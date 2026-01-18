package main

import (
	"embed"
	"io/fs"
)

//go:embed frontend/dist
var FrontendFS embed.FS

func GetFrontendFS() fs.FS {
	return FrontendFS
}
