package fsm

import (
	"testing"
	"truco/pkg/ar"
)

func TestTrucoFlow(t *testing.T) {
	m := NewMatch()

	// Initial state: PlayingState (id 1)
	if m.stateId() != 1 {
		t.Errorf("expected state 1, got %d", m.stateId())
	}

	err := m.ask(0)
	if err != nil {
		t.Fatalf("failed to ask truco: %v", err)
	}

	// State should be RespondingState (id 3)
	if m.stateId() != 3 {
		t.Errorf("expected state 3, got %d", m.stateId())
	}
	if m.cPlayer != 0 {
		t.Errorf("expected cPlayer 0 (responding should not change cPlayer), got %d", m.cPlayer)
	}

	// Player 1 accepts
	err = m.accept()
	if err != nil {
		t.Fatalf("failed to accept truco: %v", err)
	}

	// State should be back to PlayingState
	if m.stateId() != 1 {
		t.Errorf("expected state 1, got %d", m.stateId())
	}
	if m.cTruco != 2 {
		t.Errorf("expected cTruco 2, got %d", m.cTruco)
	}
}

func TestEnvidoFlow(t *testing.T) {
	m := NewMatch()

	m.cPlayer = 2
	err := m.ask(1) // envido
	if err != nil {
		t.Fatalf("failed to ask envido: %v", err)
	}

	if m.stateId() != 3 {
		t.Errorf("expected state 3, got %d", m.stateId())
	}
	if m.cPlayer != 2 {
		t.Errorf("cPlayer should not change, expected 2, got %d", m.cPlayer)
	}

	// Player 3 accepts
	err = m.accept()
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
	err = m.announce(30)
	if err != nil {
		t.Fatalf("failed to announce 30: %v", err)
	}

	// Next is P1
	err = m.announce(25)
	if err != nil {
		t.Fatalf("failed to announce 25: %v", err)
	}
	if m.envidos[1] != 130 {
		t.Fatalf("expected envido to be 130: '30 son buenas', got %v", err)
	}

	// P2
	err = m.announce(70)
	if err == nil || m.cPlayerE() != 2 {
		t.Fatalf("should fail to announce wrong envido, but didnt")
	}
	err = m.announce(13)
	if err == nil || m.cPlayerE() != 2 {
		t.Fatalf("should fail to announce wrong envido, but didnt")
	}

	m.fold()
	if m.envidos[2] != 130 {
		t.Fatalf("expected envido to be 130: '30 son buenas', got %v", err)
	}

	// P3 also folds
	m.fold()

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
	m.ask(0)
	m.fold() // P1 folds

	if m.stateId() != 0 { // EndState
		t.Errorf("expected end state, got %d", m.stateId())
	}
	if m.winnerT != 0 { // P0 wins because P1 folded
		t.Errorf("expected winner 0, got %d", m.winnerT)
	}

	err := m.play(ar.Card{N: 1, S: 'e'})
	if err == nil {
		t.Error("expected error playing in EndState")
	}

	err = m.ask(0)
	if err == nil {
		t.Error("expected error asking in EndState")
	}

	err = m.accept()
	if err == nil {
		t.Error("expected error accepting in EndState")
	}

	err = m.announce(30)
	if err == nil {
		t.Error("expected error announcing in EndState")
	}
}
