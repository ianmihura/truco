package fsm

import (
	"fmt"
	"truco/pkg/ar"
)

type AnnouncingState struct {
	match *Match
}

func (a *AnnouncingState) play(card ar.Card) error {
	return fmt.Errorf("You must announce your envido")
}

func (a *AnnouncingState) ask(requestE uint8) error {
	return fmt.Errorf("You must announce your envido")
}

func (a *AnnouncingState) accept() error {
	return fmt.Errorf("You must announce your envido")
}

// Players announce 'son buenas' by folding
func (a *AnnouncingState) fold() {
	if a.match.isEnvidoFull() {
		// should never happen
		a.match.cState = a.match.playing
		a.match.isEnvido = false
		return

	} else {
		highestE, _ := a.match.winnerE()
		a.match.envidos[a.match.cPlayerE()] = highestE + 100
	}

	if a.match.isEnvidoFull() {
		a.match.cState = a.match.playing
		a.match.isEnvido = false
	}
}

func (a *AnnouncingState) announce(score uint8) error {
	if a.match.isEnvidoFull() {
		// should never happen
		a.match.cState = a.match.playing
		a.match.isEnvido = false
		return fmt.Errorf("Announcing already finished")
	}

	highestE, _ := a.match.winnerE()

	if score <= 7 || (score >= 20 && score <= ar.MAX_ENVIDO) {
		if highestE < score {
			a.match.envidos[a.match.cPlayerE()] = score
		} else {
			// player announced loosing envido (lower than highest): should never happen
			a.match.fold()
		}

	} else {
		return fmt.Errorf("Score must be a valid envido")
	}

	if a.match.isEnvidoFull() {
		a.match.cState = a.match.playing
		a.match.isEnvido = false
	}

	return nil
}

func (a *AnnouncingState) stateId() uint8 {
	return 2
}
