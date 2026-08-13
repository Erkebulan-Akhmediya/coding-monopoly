package main

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed all:dist
var embeddedDist embed.FS

// staticHandler serves the Vue production build embedded in the binary.
// Missing paths fall back to index.html so ?admin=1 and deep links still work.
func staticHandler() http.Handler {
	sub, err := fs.Sub(embeddedDist, "dist")
	if err != nil {
		panic("embedded dist subtree missing: " + err.Error())
	}
	fileServer := http.FileServer(http.FS(sub))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			serveIndex(w, r, sub)
			return
		}

		// Reject path traversal attempts before Stat.
		if strings.Contains(path, "..") {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		info, err := fs.Stat(sub, path)
		if err != nil || info.IsDir() {
			serveIndex(w, r, sub)
			return
		}

		// Hashed Vite assets can be cached aggressively; HTML should not.
		if strings.HasPrefix(path, "assets/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			w.Header().Set("Cache-Control", "no-cache")
		}
		fileServer.ServeHTTP(w, r)
	})
}

func serveIndex(w http.ResponseWriter, r *http.Request, sub fs.FS) {
	w.Header().Set("Cache-Control", "no-cache")
	f, err := sub.Open("index.html")
	if err != nil {
		http.Error(w, "frontend not embedded", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		http.Error(w, "frontend not embedded", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodHead {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		return
	}

	rs, ok := f.(io.ReadSeeker)
	if !ok {
		data, readErr := io.ReadAll(f)
		if readErr != nil {
			http.Error(w, "frontend not embedded", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
		return
	}
	http.ServeContent(w, r, "index.html", stat.ModTime(), rs)
}
