package fsm

import (
	"fmt"
	"truco/pkg/ar"
)

type RespondingState struct {
	match *Match
}

func (r *RespondingState) play(card ar.Card) error {
	return fmt.Errorf("You cannot play, accept first")
	// _ = r.accept()
	// return r.match.Play(card)
}

func (r *RespondingState) ask(requestE AskRequest) error {
	return fmt.Errorf("You cannot ask, accept first")
	// if !r.match.IsEnvido && requestE != RequestTruco {
	// 	_, winnerE := r.match.winnerE()
	// 	if winnerE == 0 && r.match.cTurn() == 0 {
	// 		// TODO el envido va primero
	// 		return nil
	// 	} else {
	// 		return fmt.Errorf("You cannot ask for truco")
	// 	}
	// } else {
	// 	_ = r.accept()
	// 	return r.match.Ask(requestE)
	// }
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
		r.match.WinnerT = r.match.CPlayer // TODO check if this is true always
		r.match.CState = r.match.End
	}
}

func (r *RespondingState) announce(score uint8) error {
	return fmt.Errorf("You must respond")
}

func (r *RespondingState) stateId() uint8 {
	return 3
}

func (r *RespondingState) validActions() []string {
	return []string{"accept", "fold"}
}
