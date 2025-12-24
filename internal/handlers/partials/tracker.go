package partials

import (
	"html/template"
	"net/http"
)

type TrackerHandler struct {
	tmpl *template.Template
}

func NewTrackerHandler(tmpl *template.Template) *TrackerHandler {
	return &TrackerHandler{tmpl: tmpl}
}

func (h *TrackerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Simply return the player_card partial
	err := h.tmpl.ExecuteTemplate(w, "player_card", nil)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
	}
}
