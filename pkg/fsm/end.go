package fsm

import (
	"fmt"
	"truco/pkg/ar"
)

type EndState struct {
	match *Match
}

func (e *EndState) play(card ar.Card) error {
	return fmt.Errorf("Can't play a finished game")
}

func (e *EndState) ask(requestE AskRequest) error {
	return fmt.Errorf("Can't play a finished game")
}

func (e *EndState) accept() error {
	return fmt.Errorf("Can't play a finished game")
}

func (e *EndState) fold() {}

func (e *EndState) announce(score uint8) error {
	return fmt.Errorf("Can't play a finished game")
}

func (e *EndState) stateId() uint8 {
	return 0
}

func (e *EndState) validActions() []string {
	return []string{}
}
