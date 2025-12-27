package fsm

import (
	"fmt"
	"truco/pkg/ar"
)

type RespondingState struct {
	match *Match
}

// Playing is accepting
func (r *RespondingState) play(card ar.Card) error {
	// TODO only if its my turn
	_ = r.accept()
	return r.match.play(card)
}

// Asking is accepting
func (r *RespondingState) ask(requestE uint8) error {
	// TODO el envido va primero

	_ = r.accept()
	return r.match.ask(requestE)
}

func (r *RespondingState) accept() error {
	if r.match.isEnvido {
		// TODO can i accept
		r.match.cState = r.match.announcing
	} else {
		// TODO can i accept
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
