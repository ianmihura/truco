package pages

import (
	"html/template"
	"net/http"
	"truco/internal/handlers/partials"
	"truco/pkg/fsm"
	"truco/pkg/truco"
)

type MatrixHandler struct {
	Tmpl *template.Template
}

func NewMatrixHandler(tmpl *template.Template) *MatrixHandler {
	return &MatrixHandler{Tmpl: tmpl}
}

func (h *MatrixHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Compute initial stats (default to full matrix mode)
	// stats, err := truco.ComputePairStats(true, truco.FilterHands{})
	// if err != nil {
	// 	http.Error(w, "Failed to compute stats: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// statsJSON, err := json.Marshal(stats)
	// if err != nil {
	// 	http.Error(w, "Failed to marshal stats: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// First default action
	match := fsm.NewMatch()
	trackerData := partials.TrackerData{
		ActionTitle: "Jugador 1",
		Actions:     match.ValidActions(),
		State:       string(match.Encode()),
	}

	// Initial load of all cards
	cards := partials.GetAvailableCards(truco.FilterHands{})

	data := struct {
		// Stats   template.JS
		Tracker partials.TrackerData
		Cards   []partials.CardUI
	}{
		// Stats:   template.JS(statsJSON),
		Tracker: trackerData,
		Cards:   cards,
	}

	if err := h.Tmpl.ExecuteTemplate(w, "index_matrix.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
