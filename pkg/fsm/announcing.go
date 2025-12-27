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

func (a *AnnouncingState) ask(requestE AskRequest) error {
	return fmt.Errorf("You must announce your envido")
}

func (a *AnnouncingState) accept() error {
	return fmt.Errorf("You must announce your envido")
}

// Players announce 'son buenas' by folding
func (a *AnnouncingState) fold() {
	if a.match.isEnvidoFull() {
		// should never happen
		a.match.CState = a.match.Playing
		a.match.IsEnvido = false
		return

	} else {
		highestE, _ := a.match.winnerE()
		a.match.Envidos[a.match.cPlayerE()] = highestE + 100
	}

	if a.match.isEnvidoFull() {
		a.match.CState = a.match.Playing
		a.match.IsEnvido = false
	}
}

func (a *AnnouncingState) announce(score uint8) error {
	if a.match.isEnvidoFull() {
		// should never happen
		a.match.CState = a.match.Playing
		a.match.IsEnvido = false
		return fmt.Errorf("Announcing already finished")
	}

	highestE, _ := a.match.winnerE()

	if score <= 7 || (score >= 20 && score <= ar.MAX_ENVIDO) {
		if highestE < score {
			a.match.Envidos[a.match.cPlayerE()] = score
		} else {
			// player announced loosing envido (lower than highest): should never happen
			a.match.Fold()
		}

	} else {
		return fmt.Errorf("Score must be a valid envido")
	}

	if a.match.isEnvidoFull() {
		a.match.CState = a.match.Playing
		a.match.IsEnvido = false
	}

	return nil
}

func (a *AnnouncingState) stateId() uint8 {
	return 2
}

func (a *AnnouncingState) validActions() []string {
	return []string{"announce", "fold"}
}
