package fsm

import (
	"truco/pkg/ar"
)

type AskRequest uint8

const (
	RequestTruco  AskRequest = 0
	RequestEnvido AskRequest = 2
	RequestReal   AskRequest = 3
	RequestFalta  AskRequest = 255
)

// FSM for a single match
type Match struct {
	// context
	cards      [][]ar.Card // list of cards played: cards[player][turn]
	cTruco     uint8       // current truco bet (1-4)
	cTrucoAsk  uint8       // who asked for the last truco bet
	cPlayer    uint8       // player that should perform next action (not for envido)
	envidos    []uint8     // list of envidos declared per player: envidos[player] (default=255)
	cEnvido    uint8       // current envido bet 'quiero'
	cEnvidoNo  uint8       // current envido bet 'no quiero'
	cEnvidoAsk uint8       // who asked for the last envido bet
	isEnvido   bool        // so we don't duplicate response actions and states: false=truco (default), true=envido
	winnerT    uint8       // id of a player in the team that won truco, 255 if still playing
	// players are indexed as the match order:
	// 	- counter-clockwise, dealer last
	//  - 255=none

	// envidos are noted as:
	//  - 0-33:    full score
	//  - 100-133: 'son buenas': winner_env + 100
	//  - 255:     undeclared

	// states
	playing    State // can play a card or ask for truco
	responding State // can respond to asked bet
	announcing State // can announce envido amount
	end        State
	cState     State // current state
}

type Score struct {
	winnerT uint8 // player id winner of truco (unfinished=0)
	pointsT uint8 // points won in envido (default=1)

	winnerE uint8 // player id winner of envido (unplayed=0, unfinished=current)
	pointsE uint8 // points won in envido (unplayed=0)
}

// A single possible state of the game:
// interface that implements all possible actions.
// Identify the state by State.id()
type State interface {
	play(ar.Card) error            // play a card
	ask(requestE AskRequest) error // ask for a bet increase (truco, or envido with size)
	accept() error                 // accepts a bet increase
	fold()                         // rejects a bet increase, 'son buenas' in envido, or simply ends match
	announce(uint8) error          // announce how much envido you have
	stateId() uint8
}

// Returns an empty object, with binding to all states
func NewMatch() *Match {
	cards := make([][]ar.Card, 4)
	for i := range cards {
		cards[i] = make([]ar.Card, 3)
	}

	envidos := make([]uint8, 4)
	for i := range envidos {
		envidos[i] = 255
	}

	m := &Match{
		cards:      cards,
		cTruco:     1,
		cTrucoAsk:  255,
		envidos:    envidos,
		cEnvido:    0,
		cEnvidoNo:  1,
		cEnvidoAsk: 255,
		cPlayer:    0,
		isEnvido:   false,
		winnerT:    255,
	}

	m.bindStates()
	m.cState = m.playing

	return m
}

// Binds the match to all states
func (m *Match) bindStates() {
	m.playing = &PlayingState{match: m}
	m.responding = &RespondingState{match: m}
	m.announcing = &AnnouncingState{match: m}
	m.end = &EndState{match: m}
}

// Plays a card
func (m *Match) play(card ar.Card) error {
	return m.cState.play(card)
}

// Ask for a bet increase, envido or truco
func (m *Match) ask(requestE AskRequest) error {
	return m.cState.ask(requestE)
}

// Accept a bet increase
func (m *Match) accept() error {
	return m.cState.accept()
}

// If envido: rejects a bet increase.
// If declaring envido score: 'son buenas'.
// Else: ends match.
func (m *Match) fold() {
	m.cState.fold()
}

// Announce envido score
func (m *Match) announce(score uint8) error {
	return m.cState.announce(score)
}

func (m *Match) stateId() uint8 {
	return m.cState.stateId()
}

// Truco player order
func (m *Match) prevPlayer() uint8 {
	return (m.cPlayer - 1) % 4
}

// Truco player order
func (m *Match) nextPlayer() uint8 {
	return (m.cPlayer + 1) % 4
}

// Current turn, 255=end
func (m *Match) cTurn() uint8 {
	for t := range m.cards {
		for p := range m.cards[t] {
			if m.cards[t][p].N == 0 {
				return uint8(p)
			}
		}
	}
	return 255
}

// Will return true if all players declared envido,
// false if there is at least one didn't (envidos[i] == 255).
//
// Note that it returns false if envido is never played.
// (envidos array is initialized as full 255).
func (m *Match) isEnvidoFull() bool {
	return m.cPlayerE() == 255
}

// Return index of next player that needs to declare,
// returns 255 if all players declared already
func (m *Match) cPlayerE() int {
	for i := range m.envidos {
		if m.envidos[i] == 255 {
			return i
		}
	}
	return 255
}

// Returns winner envido and player id, played until now
//
// If envido is not asked, returns (0, 0)
// If envido is 'no quiero', returns (0, score)
func (m *Match) winnerE() (highest uint8, player uint8) {
	highest = 0
	if m.envidos[0] == 255 && m.cEnvido != 0 {
		// envido asked, response was 'no quiero'
		return highest, m.cEnvidoAsk
	}

	for i := range m.envidos {
		cEnv := m.envidos[i]
		if cEnv == 255 {
			// unfinished round
			break

		} else if cEnv > 100 {
			continue

		} else if cEnv > highest {
			highest = cEnv
			player = uint8(i)
		}
	}
	return highest, player
}

func (m *Match) getScore() *Score {
	winnerE, _ := m.winnerE()
	return &Score{
		winnerT: m.winnerT,
		pointsT: m.cTruco,
		winnerE: winnerE,
		pointsE: m.cEnvido,
	}
}
