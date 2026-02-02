package pages

import (
	"encoding/json"
	"html/template"
	"net/http"
	"truco/internal/handlers/partials"
	"truco/pkg/fsm"
	"truco/pkg/truco"
)

type HomeHandler struct {
	Tmpl *template.Template
}

func NewHomeHandler(tmpl *template.Template) *HomeHandler {
	return &HomeHandler{Tmpl: tmpl}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Compute initial stats (default to full matrix mode)
	stats, err := truco.ComputePairStats(true, truco.FilterHands{})
	if err != nil {
		http.Error(w, "Failed to compute stats: "+err.Error(), http.StatusInternalServerError)
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
		// Cards   []partials.CardUI
	}{
		Stats:   template.JS(statsJSON),
		Tracker: trackerData,
		// Cards:   partials.GetAvailableCards(match.GetStatsFilter()),
	}

	if err := h.Tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
