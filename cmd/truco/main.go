package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
)

// Global state for the counter (simple demo)
var (
	counter int
	mu      sync.Mutex
)

func main() {
	// Parse templates
	tmpl, err := template.ParseGlob("web/template/*.html")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	// Handlers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.ExecuteTemplate(w, "index.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// payload, err := io.ReadAll(r.Body)
		// if err != nil {
		// 	http.Error(w, "Error reading request body", http.StatusInternalServerError)
		// 	return
		// }

		// fmt.Println(payload)
		// fmt.Println(r.Body)

		mu.Lock()
		counter++
		currentVal := counter
		mu.Unlock()

		// Return just the new number as HTML (HTMX swaps this into the target)
		fmt.Fprintf(w, "%d", currentVal)
	})

	// Start server
	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
