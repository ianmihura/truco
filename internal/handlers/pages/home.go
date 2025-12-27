package pages

import (
	"encoding/json"
	"html/template"
	"net/http"
	"truco/internal/handlers/partials"
	"truco/pkg/fsm"
)

type HomeHandler struct {
	Tmpl *template.Template
	Data map[string]float64
}

func NewHomeHandler(tmpl *template.Template) *HomeHandler {
	// Load stats (hardcoded path for now, typical in this setup)
	// In a real app, inject this dependency or config.

	// data, err := ar.PairStrengths("truco_strength.csv")
	// if err != nil {
	// 	// Log error but continue with empty data to avoid crash
	// 	fmt.Printf("Error loading stats: %v\n", err)
	// 	data = make(map[string]float64)
	// }
	return &HomeHandler{Tmpl: tmpl}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jsonData, err := json.Marshal(h.Data)
	if err != nil {
		http.Error(w, "Failed to marshal stats", http.StatusInternalServerError)
		return
	}

	match := fsm.NewMatch()
	trackerData := partials.TrackerData{
		PlayerName: "Jugador 1",
		Actions:    match.ValidActions(),
		State:      string(match.Encode()),
	}

	data := struct {
		Stats   template.JS
		Tracker partials.TrackerData
	}{
		Stats:   template.JS(jsonData),
		Tracker: trackerData,
	}

	if err := h.Tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
