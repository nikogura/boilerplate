package ui

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

//go:embed all:static
var content embed.FS

// Handler returns a handler that serves static UI files from the embedded filesystem.
func Handler() (handler http.HandlerFunc) {
	// Get the filesystem with the static subdirectory
	fsys, err := fs.Sub(content, "static")
	if err != nil {
		panic("failed to create UI filesystem: " + err.Error())
	}

	fileServer := http.FileServer(http.FS(fsys))

	handler = func(w http.ResponseWriter, r *http.Request) {
		// Handle SPA routing - serve index.html for non-asset routes
		if !strings.Contains(r.URL.Path, ".") && r.URL.Path != "/" {
			// Check if the requested path exists as a file
			_, statErr := fs.Stat(fsys, strings.TrimPrefix(r.URL.Path, "/"))
			if statErr != nil {
				// File doesn't exist, serve index.html for SPA routing
				r.URL.Path = "/"
			}
		}

		// Set appropriate headers for SPA
		if strings.HasSuffix(r.URL.Path, ".html") || r.URL.Path == "/" {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
		} else {
			// Cache static assets for 1 hour
			w.Header().Set("Cache-Control", "public, max-age=3600")
		}

		// Serve the file
		fileServer.ServeHTTP(w, r)
	}

	return handler
}

// RegisterRoutes adds UI routes to the given mux router.
func RegisterRoutes(router *mux.Router) {
	// Serve static assets (CSS, JS, images)
	router.PathPrefix("/_next/").Handler(Handler())
	router.PathPrefix("/images/").Handler(Handler())
	router.PathPrefix("/static/").Handler(Handler())

	// Serve the index.html for the root and any SPA routes
	router.HandleFunc("/", Handler())
	router.NotFoundHandler = Handler()
}
