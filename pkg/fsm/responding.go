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
	// return r.match.play(card)
}

func (r *RespondingState) ask(requestE AskRequest) error {
	return fmt.Errorf("You cannot ask, accept first")
	// if !r.match.isEnvido && requestE != RequestTruco {
	// 	_, winnerE := r.match.winnerE()
	// 	if winnerE == 0 && r.match.cTurn() == 0 {
	// 		// TODO el envido va primero
	// 		return nil
	// 	} else {
	// 		return fmt.Errorf("You cannot ask for truco")
	// 	}
	// } else {
	// 	_ = r.accept()
	// 	return r.match.ask(requestE)
	// }
}

func (r *RespondingState) accept() error {
	if r.match.isEnvido {
		r.match.cState = r.match.announcing
	} else {
		r.match.cTruco += 1
		r.match.cState = r.match.playing
	}

	return nil
}

func (r *RespondingState) fold() {
	if r.match.isEnvido {
		r.match.isEnvido = false
	} else {
		r.match.winnerT = r.match.cPlayer // TODO check if this is true always
		r.match.cState = r.match.end
	}
}

func (r *RespondingState) announce(score uint8) error {
	return fmt.Errorf("You must respond")
}

func (r *RespondingState) stateId() uint8 {
	return 3
}
