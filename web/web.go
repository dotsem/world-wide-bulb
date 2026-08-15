// Package web exposes the static build output from SvelteKit for Go embed.
package web

import (
	"embed"
	"io/fs"
)

// Assets holds the static build files for the Svelte frontend.
//
//go:embed all:build
var Assets embed.FS

// GetFS returns an fs.FS sub-filesystem rooted inside "build".
func GetFS() (fs.FS, error) {
	return fs.Sub(Assets, "build")
}
