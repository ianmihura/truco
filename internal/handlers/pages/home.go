package pages

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"truco/pkg/ar"
)

type HomeHandler struct {
	Tmpl *template.Template
	Data map[string]float64
}

func NewHomeHandler(tmpl *template.Template) *HomeHandler {
	// Load stats (hardcoded path for now, typical in this setup)
	// In a real app, inject this dependency or config.
	data, err := ar.PairStrengths("truco_strength.csv")
	if err != nil {
		// Log error but continue with empty data to avoid crash
		fmt.Printf("Error loading stats: %v\n", err)
		data = make(map[string]float64)
	}
	return &HomeHandler{Tmpl: tmpl, Data: data}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jsonData, err := json.Marshal(h.Data)
	if err != nil {
		http.Error(w, "Failed to marshal stats", http.StatusInternalServerError)
		return
	}
	// passing template.JS allows the content to be written unescaped into the script tag
	if err := h.Tmpl.ExecuteTemplate(w, "index.html", template.JS(jsonData)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
