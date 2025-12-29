package pages

import (
	"encoding/json"
	"html/template"
	"net/http"
	"truco/internal/handlers/partials"
	"truco/pkg/ar"
	"truco/pkg/fsm"
)

type HomeHandler struct {
	Tmpl *template.Template
}

func NewHomeHandler(tmpl *template.Template) *HomeHandler {
	return &HomeHandler{Tmpl: tmpl}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	stats, err := ar.LoadPairStats("web/static/pair_stats.csv")
	if err != nil {
		http.Error(w, "Failed to load stats: "+err.Error(), http.StatusInternalServerError)
		return
	}

	statsJSON, err := json.Marshal(stats)
	if err != nil {
		http.Error(w, "Failed to marshal stats: "+err.Error(), http.StatusInternalServerError)
		return
	}

	match := fsm.NewMatch()
	trackerData := partials.TrackerData{
		ActionTitle: "Jugador 1",
		Actions:     match.ValidActions(),
		State:       string(match.Encode()),
	}

	data := struct {
		Stats   template.JS
		Tracker partials.TrackerData
	}{
		Stats:   template.JS(statsJSON),
		Tracker: trackerData,
	}

	if err := h.Tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
