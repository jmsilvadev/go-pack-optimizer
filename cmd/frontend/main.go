package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func serveHTML(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("static", "index.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		return
	}

	backendAPI := os.Getenv("BACKEND_URL")
	if backendAPI == "" {
		backendAPI = "http://localhost:8080/v1"
	}

	err = tmpl.Execute(w, map[string]string{
		"BackendAPI": backendAPI,
	})
	if err != nil {
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
	}
}

func main() {
	addr := os.Getenv("FRONTEND_URL")
	if addr == "" {
		addr = "0.0.0.0:3000"
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", serveHTML)

	log.Printf("Frontend running at http://%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
