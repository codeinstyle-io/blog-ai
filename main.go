package main

import (
	"embed"
	"log"
	"net/http"
	"strings"
)

//go:embed home/*
var content embed.FS

func main() {
	// Serve embedded static files
	fs := http.FileServer(http.FS(content))

	// Serve skills.json with CORS headers
	http.HandleFunc("/skills", skillsHandler)

	// Serve root and other files
	http.HandleFunc("/", rootHandler(fs))

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func skillsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	data, err := content.ReadFile("home/static/js/skills.json")
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func rootHandler(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Redirect root to index.html
		if r.URL.Path == "/" {
			data, err := content.ReadFile("home/index.html")
			if err != nil {
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(data)
			return
		}
		// Serve embedded files for matching paths
		if strings.HasSuffix(r.URL.Path, ".html") ||
			strings.HasSuffix(r.URL.Path, ".css") ||
			strings.HasSuffix(r.URL.Path, ".js") ||
			strings.HasSuffix(r.URL.Path, ".json") {
			fs.ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}
	}
}
