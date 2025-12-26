package partials

import (
	"html/template"
	"net/http"
)

type TrackerHandler struct {
	tmpl *template.Template
	data TrackerData
}

type TrackerData struct {
	PlayerName string
	Actions    []Action
}

// A single action a player can take
type Action struct {
	Name string
}

func NewTrackerHandler(tmpl *template.Template) *TrackerHandler {
	trackerData := TrackerData{
		PlayerName: "tests",
		Actions: []Action{
			{"asd"},
			{"asd"},
		},
	}
	return &TrackerHandler{tmpl: tmpl, data: trackerData}
}

func (h *TrackerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// jsonData, err := json.Marshal(h.data)
	// if err != nil {
	// 	http.Error(w, "Failed to marshal stats", http.StatusInternalServerError)
	// 	return
	// }

	err := h.tmpl.ExecuteTemplate(w, "action", h.data)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
	}
}
