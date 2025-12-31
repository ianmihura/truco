package ar

import (
	"slices"
	"testing"
	"truco/pkg/truco"
)

// Test for isEveryCardIncluded
func TestIsEveryCardIncluded(t *testing.T) {
	// Positive case: cards are included
	cards := []truco.Card{{1, 'e'}, {7, 'e'}}
	iCards := []truco.Card{{1, 'e'}, {7, 'e'}, {3, 'b'}}
	if !isEveryCardIncluded(cards, iCards) {
		t.Errorf("Expected isEveryCardIncluded to return true, got false")
	}

	// Negative case: cards are NOT included
	cards = []truco.Card{{1, 'e'}, {2, 'e'}}
	// iCards same as above
	if isEveryCardIncluded(cards, iCards) {
		t.Errorf("Expected isEveryCardIncluded to return false, got true")
	}
}

func TestCardRangeNoEnvido(t *testing.T) {
	aCards := []truco.Card{
		{1, 'e'}, {7, 'e'}, {3, 'b'}, {4, 'c'},
	}
	mCards := []truco.Card{{1, 'e'}}

	hands := cardRangeNoEnvido(aCards, mCards)

	for _, h := range hands {
		if !slices.Contains(h, mCards[0]) {
			t.Errorf("Expected hand %v to have %v in cardRangeNoEnvido output", h, mCards[0])
		}
	}
}

func TestCardRange(t *testing.T) {
	mCards := []truco.Card{{6, 'e'}}
	kCards := []truco.Card{{1, 'b'}}

	hands := CardRange(33, mCards, kCards)

	for _, h := range hands {
		if h.Envido() != 33 {
			t.Errorf("Expected hand %v to have 33 envido", h)
		}
	}
	if len(hands) != 40 {
		t.Errorf("Expected count == 40")
	}

	hands = CardRange(27, []truco.Card{{10, 'e'}, {4, 'e'}}, []truco.Card{{10, 'o'}, {5, 'e'}})
	unexpectedHand := Hand{{10, 'e'}, {4, 'e'}, {7, 'e'}}
	for _, h := range hands {
		if equalHands(h, unexpectedHand) {
			t.Errorf("Did not expect hand %v (Envido != 27) to be present", unexpectedHand)
		}
	}
	rangeHands := CardRange(127, []truco.Card{{1, 'e'}}, []truco.Card{})
	if !slices.ContainsFunc(rangeHands, func(h Hand) bool { return h.Envido() <= 27 }) {
		t.Errorf("Expected hands to have at least 27 envido")
	}

	hands = CardRange(127, []truco.Card{{10, 'c'}, {4, 'e'}}, []truco.Card{{10, 'o'}, {5, 'e'}})
	for _, h := range hands {
		if h.Envido() > 27 {
			t.Errorf("Expected hand %v to <= 27 envido", h)
		}
	}
}
