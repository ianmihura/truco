package fsm

import (
	"testing"
	"truco/pkg/truco"
)

func TestNewMatch(t *testing.T) {
	m := NewMatch()
	if m.CTruco != 1 {
		t.Errorf("expected CTruco 1, got %d", m.CTruco)
	}
	if m.CPlayer != 0 {
		t.Errorf("expected CPlayer 0, got %d", m.CPlayer)
	}
	if m.CState != m.Playing {
		t.Errorf("expected initial state to be playing")
	}
	if len(m.Envidos) != NUM_PLAYERS {
		t.Errorf("expected NUM_PLAYERS players in envidos, got %d", len(m.Envidos))
	}
	for i, e := range m.Envidos {
		if e != 255 {
			t.Errorf("expected player %d envido to be 255, got %d", i, e)
		}
	}
}

func TestPlayerOrder(t *testing.T) {
	m := NewMatch()
	m.CPlayer = 0
	if m.nextPlayer() != 1 {
		t.Errorf("next of 0 should be 1, got %d", m.nextPlayer())
	}
	if m.prevPlayer() != 3 {
		t.Errorf("prev of 0 should be 3, got %d", m.prevPlayer())
	}

	m.CPlayer = 3
	if m.nextPlayer() != 0 {
		t.Errorf("next of 3 should be 0, got %d", m.nextPlayer())
	}
	if m.prevPlayer() != 2 {
		t.Errorf("prev of 3 should be 2, got %d", m.prevPlayer())
	}
}

func TestCTurn(t *testing.T) {
	m := NewMatch()

	// Turn 0: initially no cards played
	if m.cTurn() != 0 {
		t.Errorf("expected turn 0, got %d", m.cTurn())
	}

	// Play cards for all players in turn 0
	m.Cards[0][0] = truco.Card{N: 1, S: 'e'}
	m.Cards[1][0] = truco.Card{N: 1, S: 'b'}
	m.Cards[2][0] = truco.Card{N: 1, S: 'o'}
	m.Cards[3][0] = truco.Card{N: 1, S: 'c'}

	// Should be turn 1
	if m.cTurn() != 1 {
		t.Errorf("expected turn 1, got %d", m.cTurn())
	}

	// Play cards for all players in turn 1
	m.Cards[0][1] = truco.Card{N: 7, S: 'c'}
	m.Cards[1][1] = truco.Card{N: 7, S: 'b'}
	m.Cards[2][1] = truco.Card{N: 7, S: 'o'}
	m.Cards[3][1] = truco.Card{N: 7, S: 'e'}

	// Should be turn 2
	if m.cTurn() != 2 {
		t.Errorf("expected turn 2, got %d", m.cTurn())
	}

	// Play cards for all players in turn 2
	m.Cards[0][2] = truco.Card{N: 3, S: 'c'}
	m.Cards[1][2] = truco.Card{N: 3, S: 'b'}
	m.Cards[2][2] = truco.Card{N: 3, S: 'o'}
	m.Cards[3][2] = truco.Card{N: 3, S: 'e'}

	// Should be 255 (match end)
	if m.cTurn() != 255 {
		t.Errorf("expected turn 255, got %d", m.cTurn())
	}
}

func TestEnvidoHelpers(t *testing.T) {
	m := NewMatch()

	if m.isEnvidoFull() {
		t.Errorf("expected envido not full")
	}

	if m.cPlayerE() != 0 {
		t.Errorf("expected player 0 to declare envido, got %d", m.cPlayerE())
	}

	m.Envidos[0] = 20
	if m.cPlayerE() != 1 {
		t.Errorf("expected player 1 to declare envido, got %d", m.cPlayerE())
	}

	m.Envidos[1] = 22
	m.Envidos[2] = 122
	m.Envidos[3] = 25

	if !m.isEnvidoFull() {
		t.Errorf("expected envido full")
	}
	if m.cPlayerE() != 255 {
		t.Errorf("expected playerE 255, got %d", m.cPlayerE())
	}

	highest, player := m.winnerE()
	if highest != 25 {
		t.Errorf("expected highest 25, got %d", highest)
	}
	if player != 3 {
		t.Errorf("expected winner player 3, got %d", player)
	}
}
