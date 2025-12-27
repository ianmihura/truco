package fsm

import (
	"truco/pkg/ar"
)

type AskRequest uint8

const (
	RequestTruco   AskRequest = 0
	RequestRetruco AskRequest = 1 // TODO make sure this is necesary
	RequestVale4   AskRequest = 4 // TODO make sure this is necesary
	RequestEnvido  AskRequest = 2
	RequestReal    AskRequest = 3
	RequestFalta   AskRequest = 255
	NUM_PLAYERS               = 4
)

// FSM for a single match
type Match struct {
	// context
	Cards      [][]ar.Card `json:"cards"`        // list of cards played: cards[player][turn]
	CTruco     uint8       `json:"c_truco"`      // current truco bet (1-4)
	CTrucoAsk  uint8       `json:"c_truco_ask"`  // who asked for the last truco bet
	CPlayer    uint8       `json:"c_player"`     // player that should perform next action (not for envido)
	Envidos    []uint8     `json:"envidos"`      // list of envidos declared per player: envidos[player] (default=255)
	CEnvido    uint8       `json:"c_envido"`     // current envido bet 'quiero'
	CEnvidoNo  uint8       `json:"c_envido_no"`  // current envido bet 'no quiero'
	CEnvidoAsk uint8       `json:"c_envido_ask"` // who asked for the last envido bet
	IsEnvido   bool        `json:"is_envido"`    // so we don't duplicate response actions and states: false=truco (default), true=envido
	WinnerT    uint8       `json:"winner_t"`     // id of a player in the team that won truco, 255 if still playing
	// players are indexed as the match order:
	// 	- counter-clockwise, dealer last
	//  - 255=none

	// envidos are noted as:
	//  - 0-33:    full score
	//  - 100-133: 'son buenas': winner_env + 100
	//  - 255:     undeclared

	// states
	Playing    State `json:"-"` // can play a card or ask for truco
	Responding State `json:"-"` // can respond to asked bet
	Announcing State `json:"-"` // can announce envido amount
	End        State `json:"-"` // endstate
	CState     State `json:"-"` // current state

	CStateId uint8 `json:"c_state_id"` // Helper for marshaling
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
	validActions() []string
}

// Returns an empty object, with binding to all states
func NewMatch() *Match {
	cards := make([][]ar.Card, NUM_PLAYERS)
	for i := range cards {
		cards[i] = make([]ar.Card, 3)
	}

	envidos := make([]uint8, NUM_PLAYERS)
	for i := range envidos {
		envidos[i] = 255
	}

	m := &Match{
		Cards:      cards,
		CTruco:     1,
		CTrucoAsk:  255,
		Envidos:    envidos,
		CEnvido:    0,
		CEnvidoNo:  1,
		CEnvidoAsk: 255,
		CPlayer:    0,
		IsEnvido:   false,
		WinnerT:    255,
	}

	m.bindStates()
	m.CState = m.Playing
	m.CStateId = m.CState.stateId()

	return m
}

// Binds the match to all states
func (m *Match) bindStates() {
	m.Playing = &PlayingState{match: m}
	m.Responding = &RespondingState{match: m}
	m.Announcing = &AnnouncingState{match: m}
	m.End = &EndState{match: m}
}

// Plays a card
func (m *Match) Play(card ar.Card) error {
	return m.CState.play(card)
}

// Ask for a bet increase, envido or truco
func (m *Match) Ask(requestE AskRequest) error {
	return m.CState.ask(requestE)
}

// Accept a bet increase
func (m *Match) Accept() error {
	return m.CState.accept()
}

// If envido: rejects a bet increase.
// If declaring envido score: 'son buenas'.
// Else: ends match.
func (m *Match) Fold() {
	m.CState.fold()
}

// Announce envido score
func (m *Match) Announce(score uint8) error {
	return m.CState.announce(score)
}

func (m *Match) stateId() uint8 {
	return m.CState.stateId()
}

func (m *Match) ValidActions() []string {
	return m.CState.validActions()
}

// Truco player order
func (m *Match) prevPlayer() uint8 {
	return (m.CPlayer - 1) % NUM_PLAYERS
}

// Truco player order
func (m *Match) nextPlayer() uint8 {
	return (m.CPlayer + 1) % NUM_PLAYERS
}

// Current turn, 255=end
func (m *Match) cTurn() uint8 {
	for t := range m.Cards[NUM_PLAYERS-1] {
		if m.Cards[NUM_PLAYERS-1][t].N == 0 {
			return uint8(t)
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
	for i := range m.Envidos {
		if m.Envidos[i] == 255 {
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
	if m.Envidos[0] == 255 && m.CEnvido != 0 {
		// envido asked, response was 'no quiero'
		return highest, m.CEnvidoAsk
	}

	for i := range m.Envidos {
		cEnv := m.Envidos[i]
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
		winnerT: m.WinnerT,
		pointsT: m.CTruco,
		winnerE: winnerE,
		pointsE: m.CEnvido,
	}
}
