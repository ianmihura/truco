package partials

import (
	"encoding/json"
	"html/template"
	"net/http"
	"truco/pkg/ar"
	"truco/pkg/fsm"
)

type TrackerHandler struct {
	tmpl *template.Template
}

type TrackerData struct {
	ActionTitle string
	Actions     []string
	State       string
	PlayedCard  string
	Stats       template.JS
}

func NewTrackerHandler(tmpl *template.Template) *TrackerHandler {
	return &TrackerHandler{tmpl: tmpl}
}

func (h *TrackerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	stateParam := r.URL.Query().Get("state")
	actionParam := r.URL.Query().Get("action")

	var match *fsm.Match
	if stateParam == "" {
		match = fsm.NewMatch()
	} else {
		match = fsm.Decode([]byte(stateParam))
	}

	switch actionParam {
	case "play":
		if h.handlePlay(w, r, match) {
			return
		}
	case "ask_truco", "ask_retruco", "ask_vale_4":
		_ = match.Ask(fsm.RequestTruco)
	case "ask_envido":
		_ = match.Ask(fsm.RequestEnvido)
	case "accept":
		_ = match.Accept()
	case "fold":
		match.Fold()
	case "announce":
		_ = match.Announce(20)
	}

	// TODO recalculate stats
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

	data := TrackerData{
		ActionTitle: "Jugador " + string(rune('1'+match.CPlayer)),
		Actions:     match.ValidActions(),
		State:       string(match.Encode()),
		Stats:       template.JS(statsJSON),
	}

	err = h.tmpl.ExecuteTemplate(w, "action", data)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
	}
}

// Returns true if we send cards.html list to the frontend: early return
func (h *TrackerHandler) handlePlay(w http.ResponseWriter, r *http.Request, match *fsm.Match) bool {
	card := r.URL.Query().Get("card")
	if card != "" {
		_ = match.Play(ar.NewCard(card))
		return false
	}

	// If no card specified, return the cards template
	cards := ar.ALL_CARDS // TODO reduce card options
	data := struct {
		Cards []ar.Card
		State string
	}{
		Cards: cards,
		State: string(match.Encode()),
	}
	err := h.tmpl.ExecuteTemplate(w, "cards", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return true
}
