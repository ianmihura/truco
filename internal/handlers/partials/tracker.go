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
	Actions     []fsm.ValidAction
	DoneActions []fsm.ValidAction
	State       string
	PlayedCard  string
	Stats       template.JS
	// Score       fsm.Score
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

	var doneActions []fsm.ValidAction

	action := fsm.ValidAction(actionParam)
	switch action {
	case fsm.PLAY:
		if h.handlePlay(w, r, match) {
			return // early return to skip return of new tracker
		}
	case fsm.ASK_T, fsm.ASK_RT, fsm.ASK_V4:
		_ = match.Ask(fsm.RequestTruco)
		// Return modal
		data := struct {
			Player uint8
			Action fsm.ValidAction
			State  string
		}{
			Player: (match.CPlayer+1)%fsm.NUM_PLAYERS + 1, // Next player to act
			Action: action,
			State:  string(match.Encode()),
		}
		err := h.tmpl.ExecuteTemplate(w, "truco_modal", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	case fsm.ASK_E, fsm.ASK_RE, fsm.ASK_FE:
		// TODO other envido types
		_ = match.Ask(fsm.RequestEnvido)
	case fsm.ACCEPT:
		prevTruco := match.CTruco
		_ = match.Accept()
		if match.CTruco > prevTruco {
			switch match.CTruco {
			case 2:
				doneActions = append(doneActions, fsm.ASK_T)
			case 3:
				doneActions = append(doneActions, fsm.ASK_RT)
			case 4:
				doneActions = append(doneActions, fsm.ASK_V4)
			}
		}
	case fsm.FOLD, fsm.FOLD_NQ, fsm.FOLD_SB:
		match.Fold()
	case fsm.ANNOUN:
		_ = match.Announce(20)
	}

	data := TrackerData{
		ActionTitle: "Jugador " + string(rune('1'+match.CPlayer)),
		Actions:     match.ValidActions(),
		DoneActions: doneActions,
		State:       string(match.Encode()),
		// Score:       *match.GetScore(),
	}

	err := h.tmpl.ExecuteTemplate(w, "action", data)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *TrackerHandler) HandleStats(w http.ResponseWriter, r *http.Request) {
	stateParam := r.URL.Query().Get("state")
	fmatrixParam := r.URL.Query().Get("fmatrix")

	var match *fsm.Match
	if stateParam == "" {
		match = fsm.NewMatch()
	} else {
		match = fsm.Decode([]byte(stateParam))
	}

	// Recalculate stats dynamically based on the current matrix mode
	stats, err := ar.ComputePairStats(fmatrixParam == "true", match.GetStatsFilter())
	if err != nil {
		http.Error(w, "Failed to compute stats: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, "Failed to marshal stats: "+err.Error(), http.StatusInternalServerError)
	}
}

// Returns true if we send cards.html list to the frontend: early return
func (h *TrackerHandler) handlePlay(w http.ResponseWriter, r *http.Request, match *fsm.Match) bool {
	card := r.URL.Query().Get("card")
	if card != "" {
		// If user specified a card, return the next action
		_ = match.Play(ar.NewCard(card))
		return false
	}
	// If no card specified, return the cards template

	cards := ar.ALL_CARDS // TODO reduce card options
	data := struct {
		Cards  []ar.Card
		State  string
		Action fsm.ValidAction
	}{
		Cards:  cards,
		State:  string(match.Encode()),
		Action: fsm.PLAY,
	}
	err := h.tmpl.ExecuteTemplate(w, "cards", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return true
}
