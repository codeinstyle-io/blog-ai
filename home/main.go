package main

import (
	"embed"
	"log"
	"net/http"
)

//go:embed static/*
var static embed.FS

//go:embed templates/*
var templates embed.FS

func skillsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	data, err := static.ReadFile("static/data/skills.json")
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func catchAllHandler(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Redirect root to index.html
		if r.URL.Path == "/" {
			data, err := templates.ReadFile("templates/index.html")

			if err != nil {
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(data)
			return
		}

		if r.URL.Path == "/about" {
			data, err := templates.ReadFile("templates/about.html")
			if err != nil {
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(data)
			return
		}

		// Serve embedded files for matching paths
		fs.ServeHTTP(w, r)
	}
}

func main() {
	// Serve embedded static files
	fs := http.FileServer(http.FS(static))

	// Serve skills.json with CORS headers
	http.HandleFunc("/skills", skillsHandler)

	// Serve root and other files
	http.HandleFunc("/", catchAllHandler(fs))

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
