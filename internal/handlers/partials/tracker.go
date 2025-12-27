package partials

import (
	"html/template"
	"net/http"
	"truco/pkg/ar"
	"truco/pkg/fsm"
)

type TrackerHandler struct {
	tmpl *template.Template
}

type TrackerData struct {
	PlayerName string
	Actions    []string
	State      string
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

	if actionParam != "" {
		switch actionParam {
		case "play":
			// Hardcoded for now as requested
			_ = match.Play(ar.Card{N: 1, S: 'e'})
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
	}

	data := TrackerData{
		PlayerName: "Jugador " + string(rune('1'+match.CPlayer)),
		Actions:    match.ValidActions(),
		State:      string(match.Encode()),
	}

	err := h.tmpl.ExecuteTemplate(w, "action", data)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
	}
}
