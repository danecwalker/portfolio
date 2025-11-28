package frontend

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

//go:generate bun i
//go:generate bun run build
//go:embed all:build
var files embed.FS

func SvelteKitHandler() http.Handler {
	// make an fs.FS rooted at build/
	static, err := fs.Sub(files, "build")
	if err != nil {
		log.Fatal(err)
	}
	fileServer := http.FileServer(http.FS(static))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// strip the leading slash for FS lookups
		upath := strings.TrimPrefix(r.URL.Path, "/")

		if upath == "" {
			// if no path, serve the index.html
			r.URL.Path = "/"
			fileServer.ServeHTTP(w, r)
			return
		}

		// if we have an extension, serve the file
		if filepath.Ext(upath) != "" {
			// serve the file
			fileServer.ServeHTTP(w, r)
			return
		}

		_, err := static.Open(upath + ".html")
		if err == nil {
			// if we have a .html file, serve it
			r.URL.Path = upath + ".html"
			fileServer.ServeHTTP(w, r)
			return
		}

		r.URL.Path = "/fallback.html"
		fileServer.ServeHTTP(w, r)
	})
}
