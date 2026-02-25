package main

import (
	"html/template"
	"log"
	"net/http"
	"truco/internal/server"
)

func main() {
	// Parse templates
	tmpl, err := template.ParseGlob("web/template/*.html")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	// Initialize Server
	srv := server.NewServer(tmpl)

	// Start server
	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", srv); err != nil {
		log.Fatal(err)
	}
}
