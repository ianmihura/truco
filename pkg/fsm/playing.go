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
		p.match.cState = p.match.end
		return p.match.play(card)
	}

	p.match.cards[p.match.cPlayer][turn] = card
	p.match.cPlayer = p.match.nextPlayer()
	return nil
}

func (p *PlayingState) ask(requestE AskRequest) error {
	if requestE != RequestTruco {
		if p.match.cTurn() == 0 {
			if !p.match.isEnvido { // first envido request
				if p.match.cPlayer >= 2 { // only last two players can request it
					p.match.cEnvidoAsk = p.match.cPlayer
					p.match.cEnvido += uint8(requestE)
					p.match.isEnvido = true
				} else {
					return fmt.Errorf("You can't ask for envido")
				}

			} else {
				p.match.cEnvidoAsk = (p.match.cEnvidoAsk + 1) % 4
				if requestE == RequestFalta {
					p.match.cEnvido = uint8(RequestFalta)
					p.match.cEnvidoNo += 1 // TODO not correct
				} else {
					p.match.cEnvido += uint8(requestE)
					p.match.cEnvidoNo += 1 // TODO not correct
				}
			}

			p.match.cState = p.match.responding
			return nil
		} else {
			return fmt.Errorf("You can't ask for envido")
		}

	} else {
		if p.match.cTruco == 4 {
			return fmt.Errorf("Truco is highest")
		}

		if p.match.cTrucoAsk%2 != p.match.cPlayer%2 {
			p.match.cTrucoAsk = p.match.cPlayer
			p.match.isEnvido = false
			// p.match.cTruco changes in accept action
			p.match.cState = p.match.responding
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
	p.match.winnerT = p.match.prevPlayer()
	p.match.cState = p.match.end
}

func (p *PlayingState) announce(score uint8) error {
	return fmt.Errorf("You must play a card or raise")
}

func (p *PlayingState) stateId() uint8 {
	return 1
}
