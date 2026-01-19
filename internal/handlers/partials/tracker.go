package partials

import (
	"net/http"
	"strconv"
	"truco/pkg/fsm"
	"truco/pkg/truco"
)

func (h *Handler) TrackAct(w http.ResponseWriter, r *http.Request) {
	actionParam := r.URL.Query().Get("action")
	match := GetMatch(r)

	action := fsm.ValidAction(actionParam)
	tmplName, data := processActionFSM(action, match, r)

	err := h.tmpl.ExecuteTemplate(w, tmplName, data)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
	}
}

func processActionFSM(action fsm.ValidAction, match *fsm.Match, r *http.Request) (string, any) {
	var doneActions []fsm.ValidAction

	switch action {
	case fsm.PLAY:
		card := r.URL.Query().Get("card")
		_ = match.Play(truco.NewCard(card))

	case fsm.ASK_T, fsm.ASK_RT, fsm.ASK_V4:
		_ = match.Ask(fsm.RequestTruco)
		return "truco_modal", struct {
			Player uint8
			Action fsm.ValidAction
			State  string
		}{
			Player: match.CPlayer + 1,
			Action: action,
			State:  string(match.Encode()),
		}

	case fsm.ASK_E, fsm.ASK_RE, fsm.ASK_FE:
		return "envido_selector", struct {
			Envidos [][]fsm.ValidAction
			Players int
			State   string
		}{
			Envidos: match.ValidEnvidos(),
			Players: fsm.NUM_PLAYERS,
			State:   string(match.Encode()),
		}

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
		// TODO revise this
		// 1. Process Betting Sequence
		idxStr := r.URL.Query().Get("combination_idx")
		if idx, err := strconv.Atoi(idxStr); err == nil {
			combos := match.ValidEnvidos()
			if idx >= 0 && idx < len(combos) {
				combo := combos[idx]
				for _, act := range combo {
					req := fsm.RequestEnvido
					switch act {
					case fsm.ASK_RE:
						req = fsm.RequestReal
					case fsm.ASK_FE:
						req = fsm.RequestFalta
					}
					_ = match.Ask(req)
				}
			}
		}
		_ = match.Accept()

		// 2. Process Announcements
		// Loop while someone needs to announce
		for {
			p := match.CPlayerE()
			if p == 255 {
				break
			}
			scoreStr := r.URL.Query().Get("score_" + strconv.Itoa(p))
			score := 0
			if s, err := strconv.Atoi(scoreStr); err == nil {
				score = s
			}
			_ = match.Announce(uint8(score))
		}
		doneActions = append(doneActions, fsm.ValidAction("Canta"))
	}

	// default: next action tracker (next player's turn)
	return "action", TrackerData{
		ActionTitle: "Jugador " + string(rune('1'+match.CPlayer)),
		Actions:     match.ValidActions(),
		DoneActions: doneActions,
		State:       string(match.Encode()),
	}
}
