package fsm

import (
	"testing"
	"truco/pkg/truco"
)

func TestTrucoFlow(t *testing.T) {
	m := NewMatch()

	// Initial state: PlayingState (id 1)
	if m.stateId() != 1 {
		t.Errorf("expected state 1, got %d", m.stateId())
	}

	err := m.Ask(RequestTruco)
	if err != nil {
		t.Fatalf("failed to ask truco: %v", err)
	}

	// State should be RespondingState (id 3)
	if m.stateId() != 3 {
		t.Errorf("expected state 3, got %d", m.stateId())
	}
	if m.CPlayer != 0 {
		t.Errorf("expected CPlayer 0 (responding should not change CPlayer), got %d", m.CPlayer)
	}

	// Player 1 accepts
	err = m.Accept()
	if err != nil {
		t.Fatalf("failed to accept truco: %v", err)
	}

	// State should be back to PlayingState
	if m.stateId() != 1 {
		t.Errorf("expected state 1, got %d", m.stateId())
	}
	if m.CTruco != 2 {
		t.Errorf("expected CTruco 2, got %d", m.CTruco)
	}
}

func TestEnvidoFlow(t *testing.T) {
	m := NewMatch()

	m.CPlayer = 2
	err := m.Ask(RequestEnvido) // envido
	if err != nil {
		t.Fatalf("failed to ask envido: %v", err)
	}

	if m.stateId() != 3 {
		t.Errorf("expected state 3, got %d", m.stateId())
	}
	if m.CPlayer != 2 {
		t.Errorf("CPlayer should not change, expected 2, got %d", m.CPlayer)
	}

	// Player 3 accepts
	err = m.Accept()
	if err != nil {
		t.Fatalf("failed to accept envido: %v", err)
	}

	// State should be AnnouncingState (id 2)
	if m.stateId() != 2 {
		t.Errorf("expected state 2, got %d", m.stateId())
	}

	// P0 starts announcing (cPlayerE returns first 255)
	if m.cPlayerE() != 0 {
		t.Errorf("expected player 0 to announce, got %d", m.cPlayerE())
	}

	// P0
	err = m.Announce(30)
	if err != nil {
		t.Fatalf("failed to announce 30: %v", err)
	}

	// Next is P1
	err = m.Announce(25)
	if err != nil {
		t.Fatalf("failed to announce 25: %v", err)
	}
	if m.Envidos[1] != 130 {
		t.Fatalf("expected envido to be 130: '30 son buenas', got %d", m.Envidos[1])
	}

	// P2
	err = m.Announce(70)
	if err == nil || m.cPlayerE() != 2 {
		t.Fatalf("should fail to announce wrong envido, but didnt")
	}
	err = m.Announce(13)
	if err == nil || m.cPlayerE() != 2 {
		t.Fatalf("should fail to announce wrong envido, but didnt")
	}

	m.Fold()
	if m.Envidos[2] != 130 {
		t.Fatalf("expected envido to be 130: '30 son buenas', got %d", m.Envidos[2])
	}

	// P3 also folds
	m.Fold()

	// Should be back to PlayingState
	if m.stateId() != 1 {
		t.Errorf("expected state 1, got %d", m.stateId())
	}

	highest, winner := m.winnerE()
	if highest != 30 || winner != 0 {
		t.Errorf("expected winner 0 with 30, got %d with %d", winner, highest)
	}
}

func TestFoldTrucoAndEndstate(t *testing.T) {
	m := NewMatch()
	m.Ask(RequestTruco)
	m.Fold() // P1 folds

	if m.stateId() != 0 { // EndState
		t.Errorf("expected end state, got %d", m.stateId())
	}
	if m.WinnerT != 0 { // P0 wins because P1 folded
		t.Errorf("expected winner 0, got %d", m.WinnerT)
	}

	err := m.Play(truco.Card{N: 1, S: 'e'})
	if err == nil {
		t.Error("expected error playing in EndState")
	}

	err = m.Ask(RequestTruco)
	if err == nil {
		t.Error("expected error asking in EndState")
	}

	err = m.Accept()
	if err == nil {
		t.Error("expected error accepting in EndState")
	}

	err = m.Announce(30)
	if err == nil {
		t.Error("expected error announcing in EndState")
	}
}
