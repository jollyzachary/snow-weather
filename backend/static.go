package main

import (
	"bytes"
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

//go:embed all:static
var embeddedStatic embed.FS

func spaHandler() http.Handler {
	staticFiles, err := fs.Sub(embeddedStatic, "static/dist")
	if err != nil {
		panic(err)
	}
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet && request.Method != http.MethodHead {
			http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		requestedPath := strings.TrimPrefix(path.Clean(request.URL.Path), "/")
		if requestedPath == "." || requestedPath == "" {
			requestedPath = "index.html"
		}

		fileInfo, err := fs.Stat(staticFiles, requestedPath)
		if err != nil || fileInfo.IsDir() {
			if path.Ext(requestedPath) != "" {
				http.NotFound(writer, request)
				return
			}
			requestedPath = "index.html"
			fileInfo, err = fs.Stat(staticFiles, requestedPath)
			if err != nil {
				http.Error(writer, "interface unavailable", http.StatusInternalServerError)
				return
			}
		}

		contents, err := fs.ReadFile(staticFiles, requestedPath)
		if err != nil {
			http.Error(writer, "interface unavailable", http.StatusInternalServerError)
			return
		}
		http.ServeContent(writer, request, requestedPath, fileInfo.ModTime(), bytes.NewReader(contents))
	})
}
