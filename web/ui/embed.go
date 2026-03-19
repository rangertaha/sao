package ui

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

//go:embed dist/*
var dist embed.FS

// FileServer returns a handler serving embedded UI assets with SPA fallback.
func FileServer() (http.Handler, error) {
	sub, err := fs.Sub(dist, "dist")
	if err != nil {
		return nil, err
	}

	fileServer := http.FileServer(http.FS(sub))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cleanPath := path.Clean(r.URL.Path)
		if cleanPath == "." {
			cleanPath = "/"
		}

		// For SPA routes, serve index.html unless path looks like an asset.
		if !strings.Contains(path.Base(cleanPath), ".") {
			r = r.Clone(r.Context())
			r.URL.Path = "/index.html"
		}

		fileServer.ServeHTTP(w, r)
	}), nil
}
