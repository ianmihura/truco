package ar

import (
	"slices"
	"testing"
)

func TestSortCards(t *testing.T) {
	strs := Hand{
		Card{12, 0},
		Card{5, 0},
		Card{7, 0},
		Card{1, 0},
		Card{10, 0},
		Card{3, 0},
		Card{11, 0},
		Card{1, 0},
	}
	slices.SortFunc(strs, SortForEnvido)

	sorted := []uint8{7, 5, 3, 1, 1, 12, 10, 11}

	for i := range len(strs) {
		if strs[i].N != sorted[i] {
			t.Errorf("Error sorting Hand, at index %d got %d, expected %d", i, strs[i].N, sorted[i])
		}
	}
}

func TestEnvidoScore(t *testing.T) {
	var envido uint8
	envido = Hand{Card{1, 'e'}, Card{7, 'e'}, Card{6, 'o'}}.Envido()
	if envido != 28 {
		t.Errorf("Error in envido got %d expected %d", envido, 28)
	}

	envido = Hand{Card{1, 'e'}, Card{7, 'e'}, Card{10, 'e'}}.Envido()
	if envido != 28 {
		t.Errorf("Error in envido got %d expected %d", envido, 28)
	}

	envido = Hand{Card{1, 'e'}, Card{7, 'e'}, Card{2, 'e'}}.Envido()
	if envido != 29 {
		t.Errorf("Error in envido got %d expected %d", envido, 29)
	}

	envido = Hand{Card{1, 'e'}, Card{7, 'b'}, Card{6, 'e'}}.Envido()
	if envido != 27 {
		t.Errorf("Error in envido got %d expected %d", envido, 27)
	}

	envido = Hand{Card{1, 'e'}, Card{7, 'b'}, Card{6, 'o'}}.Envido()
	if envido != 7 {
		t.Errorf("Error in envido got %d expected %d", envido, 7)
	}

	envido = Hand{Card{1, 'e'}, Card{7, 'o'}, Card{6, 'o'}}.Envido()
	if envido != 33 {
		t.Errorf("Error in envido got %d expected %d", envido, 33)
	}

	envido = Hand{Card{1, 'e'}, Card{11, 'o'}, Card{10, 'o'}}.Envido()
	if envido != 20 {
		t.Errorf("Error in envido got %d expected %d", envido, 20)
	}

	envido = Hand{Card{1, 'e'}, Card{11, 'e'}, Card{10, 'o'}}.Envido()
	if envido != 21 {
		t.Errorf("Error in envido got %d expected %d", envido, 21)
	}

	envido = Hand{Card{1, 'e'}, Card{11, 'o'}, Card{10, 'c'}}.Envido()
	if envido != 1 {
		t.Errorf("Error in envido got %d expected %d", envido, 1)
	}

	envido = Hand{Card{12, 'e'}, Card{11, 'o'}, Card{10, 'c'}}.Envido()
	if envido != 0 {
		t.Errorf("Error in envido got %d expected %d", envido, 0)
	}
}
func TestEnvidoPairsCount(t *testing.T) {
	for i := range 8 {
		score := uint8(i)
		if int(envidoPairsCount(score)) != len(envidoPairs(score)) {
			t.Errorf("Error in envidoPairsCount for %d: got %d, expected %d", score, envidoPairsCount(score), len(envidoPairs(score)))
		}
	}

	for i := 20; i < 34; i++ {
		score := uint8(i)
		if int(envidoPairsCount(score)) != len(envidoPairs(score)) {
			t.Errorf("Error in envidoPairsCount for %d: got %d, expected %d", score, envidoPairsCount(score), len(envidoPairs(score)))
		}
	}
}

// for testing purposes: counts expected result of EnvidoPairs
func envidoPairsCount(score uint8) uint8 {
	cardVal := ENVIDOS[score]
	partial := uint8(len(cardVal) * 4)

	if score == 0 {
		return 12
	} else if score < 20 {
		return 4
	} else if score == 20 {
		return 12
	} else if score <= 27 {
		// if they have a figure
		return partial + 8
	} else {
		return partial
	}
}

func TestEnvidoHands(t *testing.T) {
	// Test envido 33
	// EnvidoPairs(33) returns {6s, 7s} for s in [e, b, o, c] (4 pairs)
	// For each pair, there are 4 suits * 10 cards - 10 cards of same suit = 30 cards of diff suit.
	// Total hands = 4 * 30 = 120 ?
	// Actually ALL_CARDS has 40 cards.
	// Pair has suit S. 3rd card must be diff suit. There are 30 cards of diff suit.
	// So 4 pairs * 30 cards = 120 hands.
	hands := EnvidoHands(33)
	if len(hands) != 120 {
		t.Errorf("Expected 120 hands for envido 33, got %d", len(hands))
	}

	for _, h := range envido33Check {
		found := false
		for _, generated := range hands {
			if equalHands(h, generated) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected hand %v not found in generated hands for 33", h)
		}
	}

	for _, h := range hands {
		if h.Envido() != 33 {
			t.Errorf("Generated hand %v has envido %d, expected 33", h, h.Envido())
		}
	}

	// Test envido 7 (single card 7, or pair {1,6} no {1,6} is 27. Envido 7 is just 7?)
	// Actually EnvidoPairs(7) returns {{7}} for each suit? No check data.go
	// ENVIDOS[7] = {{'7'}}
	// So 4 single cards: 7e, 7b, 7o, 7c.
	// For 7e:
	// Need 2 cards c2, c3 such that suit(c2)!=e, suit(c3)!=e, suit(c2)!=suit(c3)
	// And c2.Envido <= 7, c3.Envido <= 7. (All cards have Envido <= 7 except none? Wait Envido is rank for non-figures).
	// Actually Envido value of 7 is 7. Envido value of 10,11,12 is 0.
	// So all cards have envido value <= 7.
	// So just 3 diff suits.
	// Suits: e, b, o, c. 7e takes e.
	// Pairs of suits from {b, o, c} are {b,o}, {b,c}, {o,c} (3 pairs).
	// For each suit pair, we pick any card (10 ranks).
	// So 3 * 10 * 10 = 300 hands for 7e.
	// Total 4 * 300 = 1200 hands.
	hands7 := EnvidoHands(7)
	// Let's just check length if calculation is correct or simply check validity
	if len(hands7) != 1200 {
		t.Errorf("Expected 1200 hands for envido 7, got %d", len(hands7))
	}
	for _, h := range hands7 {
		if h.Envido() != 7 {
			t.Errorf("Generated hand %v has envido %d, expected 7", h, h.Envido())
		}
	}

	// Test Envido 20
	// Calculation:
	// EnvidoPairs(20) returns pairs of figures of the same suit.
	// Figures per suit: 10, 11, 12. Pairs: {10,11}, {10,12}, {11,12} = 3 pairs.
	// Total suits: 4. Total pairs: 4 * 3 = 12.
	// For each pair, 3rd card must be diff suit. 30 options.
	// Total hands: 12 * 30 = 360.
	hands20 := EnvidoHands(20)
	if len(hands20) != 360 {
		t.Errorf("Expected 360 hands for envido 20, got %d", len(hands20))
	}
	for _, h := range hands20 {
		if h.Envido() != 20 {
			t.Errorf("Generated hand %v has envido %d, expected 20", h, h.Envido())
		}
	}

	// Test Envido 25
	// Calculation:
	// ENVIDOS[25] = {{'f', 5}, {1, 4}, {2, 3}}
	// 1. {Figure, 5}: Figures {10,11,12}. Pairs {10,5}, {11,5}, {12,5}. 3 per suit -> 12 total.
	// 2. {1, 4}: Pairs {1,4}. 1 per suit -> 4 total.
	// 3. {2, 3}: Pairs {2,3}. 1 per suit -> 4 total.
	// Total pairs = 12 + 4 + 4 = 20.
	// 3rd card diff suit: 30 options.
	// Total hands: 20 * 30 = 600.
	hands25 := EnvidoHands(25)
	if len(hands25) != 600 {
		t.Errorf("Expected 600 hands for envido 25, got %d", len(hands25))
	}
	for _, h := range hands25 {
		if h.Envido() != 25 {
			t.Errorf("Generated hand %v has envido %d, expected 25", h, h.Envido())
		}
	}
}

// Helper to check if two hands are equal (ignoring order)
func equalHands(h1, h2 Hand) bool {
	if len(h1) != len(h2) {
		return false
	}
	s1 := make([]string, len(h1))
	s2 := make([]string, len(h2))
	for i := range h1 {
		s1[i] = string(h1[i].S) + string(rune(h1[i].N)) // Rough key
	}
	for i := range h2 {
		s2[i] = string(h2[i].S) + string(rune(h2[i].N))
	}
	// Better: sort cards and compare
	// But Hand struct doesn't strictly define sort.
	// Let's just use double loop match
	used := make([]bool, len(h2))
	for _, c1 := range h1 {
		found := false
		for j, c2 := range h2 {
			if !used[j] && c1 == c2 {
				used[j] = true
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

var envido33Check = []Hand{
	{Card{6, 'e'}, Card{7, 'e'}, Card{1, 'b'}},
	{Card{6, 'e'}, Card{7, 'e'}, Card{10, 'o'}},
}
