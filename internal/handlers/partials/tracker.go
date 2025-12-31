package partials

import (
	"net/http"
	"slices"
	"truco/pkg/ar"
	"truco/pkg/fsm"
	"truco/pkg/truco"
)

func (h *Handler) TrackAct(w http.ResponseWriter, r *http.Request) {
	actionParam := r.URL.Query().Get("action")
	match := GetMatch(r)

	var doneActions []fsm.ValidAction

	action := fsm.ValidAction(actionParam)
	switch action {
	case fsm.PLAY:
		if h.handlePlay(w, r, match) {
			return // early return: skip new tracker
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
		return // early return: skip new tracker
	case fsm.ASK_E:
		_ = match.Ask(fsm.RequestEnvido)
	case fsm.ASK_RE:
		_ = match.Ask(fsm.RequestReal)
	case fsm.ASK_FE:
		_ = match.Ask(fsm.RequestFalta)
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
	}

	err := h.tmpl.ExecuteTemplate(w, "action", data)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
	}
}

// Returns true if we send cards.html list to the frontend: early return
func (h *Handler) handlePlay(w http.ResponseWriter, r *http.Request, match *fsm.Match) bool {
	card := r.URL.Query().Get("card")
	if card != "" {
		// If user specified a card, return the next action
		_ = match.Play(truco.NewCard(card))
		return false
	}
	// If no card specified, return the cards template

	slices.SortFunc(truco.ALL_CARDS, ar.SortForTruco)
	data := struct {
		Cards  []CardUI
		State  string
		Action fsm.ValidAction
	}{
		Cards:  GetAvailableCards(match.GetStatsFilter()),
		State:  string(match.Encode()),
		Action: fsm.PLAY,
	}
	err := h.tmpl.ExecuteTemplate(w, "cards", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return true
}
