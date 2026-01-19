package fsm

import (
	"fmt"
	"truco/pkg/truco"
)

type RespondingState struct {
	match *Match
}

func (r *RespondingState) play(card truco.Card) error {
	return fmt.Errorf("You cannot play, accept first")
	// _ = r.accept()
	// return r.match.Play(card)
}

func (r *RespondingState) ask(requestE AskRequest) error {
	// TODO allow "el envido va primero"

	if r.match.IsEnvido && requestE != RequestTruco {
		// Envido re-raise
		r.match.CEnvidoAsk = r.match.CPlayer
		if requestE == RequestFalta {
			r.match.CEnvido = uint8(RequestFalta)
		} else {
			r.match.CEnvido += uint8(requestE)
		}
		// TODO stay in Responding state, but now it's the other team's turn to respond
		return nil
	}
	return fmt.Errorf("You cannot ask, accept first")
}

func (r *RespondingState) accept() error {
	if r.match.IsEnvido {
		r.match.CState = r.match.Announcing
	} else {
		r.match.CTruco += 1
		r.match.CState = r.match.Playing
	}

	return nil
}

func (r *RespondingState) fold() {
	if r.match.IsEnvido {
		r.match.IsEnvido = false
		r.match.CState = r.match.Playing
	} else {
		r.match.CState = r.match.End
		r.match.WinnerT = r.match.CPlayer
		// NOTE: this works only if we keep atomic ask-response:
		// if we allow classic ask-ask-respond, it will not.
	}
}

func (r *RespondingState) announce(score uint8) error {
	return fmt.Errorf("You must respond")
}

func (r *RespondingState) stateId() uint8 {
	return 3
}

func (r *RespondingState) validActions() []ValidAction {
	actions := []ValidAction{ACCEPT, FOLD_NQ}
	if r.match.IsEnvido {
		if r.match.CEnvido < 255 {
			if r.match.CEnvido < 3 {
				actions = append(actions, ASK_E)
			}
			actions = append(actions, ASK_RE, ASK_FE)
		}
	}
	return actions
}
