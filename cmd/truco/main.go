package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"truco/internal/server"
	"truco/pkg/truco"
)

func main() {
	tmpl := template.New("").Funcs(template.FuncMap{
		"f32": func(a int) float32 {
			return float32(a)
		},
		"mul": func(a, b float32) float32 {
			return a * b
		},
		"mul_i": func(a, b float32) int {
			return int(a * b)
		},
		"sub": func(a uint8, b int) int {
			return int(a) - b
		},
		"mapCardEmoji": func(card string) string {
			return truco.NewCard(card).ToEmoji()
		},
		"thousand_int": func(n int) string {
			s := strconv.FormatInt(int64(n), 10)
			if len(s) <= 3 {
				return s
			}
			var buf strings.Builder
			for i, r := 0, len(s)%3; i < len(s); i++ {
				if i > 0 && i%3 == r {
					buf.WriteByte(',')
				}
				buf.WriteByte(s[i])
			}
			return buf.String()
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
