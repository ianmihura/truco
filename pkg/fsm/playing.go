package fsm

import (
	"fmt"
	"truco/pkg/ar"
)

type PlayingState struct {
	match *Match
}

func (p *PlayingState) play(card ar.Card) error {
	turn := p.match.cTurn()
	if turn == 255 {
		// finished match
		p.match.CState = p.match.End
		return p.match.Play(card)
	}

	p.match.Cards[p.match.CPlayer][turn] = card
	p.match.CPlayer = p.match.nextPlayer()
	return nil
}

func (p *PlayingState) ask(requestE AskRequest) error {
	if requestE != RequestTruco {
		if p.match.cTurn() == 0 {
			if !p.match.IsEnvido { // first envido request
				if p.match.CPlayer >= 2 { // only last two players can request it
					p.match.CEnvidoAsk = p.match.CPlayer
					p.match.CEnvido += uint8(requestE)
					p.match.IsEnvido = true
				} else {
					return fmt.Errorf("You can't ask for envido")
				}

			} else {
				p.match.CEnvidoAsk = (p.match.CEnvidoAsk + 1) % NUM_PLAYERS
				if requestE == RequestFalta {
					p.match.CEnvido = uint8(RequestFalta)
					p.match.CEnvidoNo += 1 // TODO not correct
				} else {
					p.match.CEnvido += uint8(requestE)
					p.match.CEnvidoNo += 1 // TODO not correct
				}
			}

			p.match.CState = p.match.Responding
			return nil
		} else {
			return fmt.Errorf("You can't ask for envido")
		}

	} else {
		if p.match.CTruco == 4 {
			return fmt.Errorf("Truco is highest")
		}

		if p.match.CTrucoAsk%2 != p.match.CPlayer%2 {
			p.match.CTrucoAsk = p.match.CPlayer
			p.match.IsEnvido = false
			// p.match.cTruco changes in accept action
			p.match.CState = p.match.Responding
			return nil
		} else {
			return fmt.Errorf("You can't ask for truco")
		}
	}
}

func (p *PlayingState) accept() error {
	return fmt.Errorf("You must play a card or raise")
}

func (p *PlayingState) fold() {
	p.match.WinnerT = p.match.prevPlayer()
	p.match.CState = p.match.End
}

func (p *PlayingState) announce(score uint8) error {
	return fmt.Errorf("You must play a card or raise")
}

func (p *PlayingState) stateId() uint8 {
	return 1
}

func (p *PlayingState) validActions() []string {
	// return []string{"play", "fold", "ask"}

	// TODO check this:
	actions := []string{"play", "fold"}
	if p.match.CTruco < 4 && p.match.CTrucoAsk%2 != p.match.CPlayer%2 {
		actions = append(actions, "ask_truco")
	}
	if p.match.cTurn() == 0 {
		if !p.match.IsEnvido {
			if p.match.CPlayer >= 2 {
				actions = append(actions, "ask_envido")
			}
		} else {
			actions = append(actions, "ask_envido")
		}
	}
	return actions
}
