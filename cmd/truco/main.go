package main

import (
	"html/template"
	"log"
	"net/http"
	"truco/internal/server"
)

func main() {
	tmpl := template.New("").Funcs(template.FuncMap{
		"mul": func(a, b float32) float32 {
			return a * b
		},
	})
	tmpl, err := tmpl.ParseGlob("web/template/*.html")
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
