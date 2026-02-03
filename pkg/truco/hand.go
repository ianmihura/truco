package truco

import (
	"fmt"
	"slices"
	"strings"
)

// Slice of Cards of any length
type Hand []Card

// Returns a new Hand from a string: Hand{NewCard(c[0]), ...}.
// Expects cards to be separated by a single space.
// Will use the function `NewCard` to create cards from each card
// (see that function for more details)
//
//	"1e 1b 1c" -> Hand{{1, 'e'}, {1, 'b'}, {1, 'c'}}
func NewHand(handStr string) Hand {
	cards := strings.Split(handStr, " ")
	hand := make([]Card, 0, len(cards))
	for c := range cards {
		hand = append(hand, NewCard(cards[c]))
	}
	return Hand(hand)
}

// Returns true if the hand has all specified cards
func (hand Hand) HasAll(kCards []Card) bool {
	if len(kCards) > len(hand) {
		return false
	}

outer:
	for _, k := range kCards {
		for _, h := range hand {
			if h == k {
				continue outer
			}
		}
		return false
	}
	return true
}

// Returns true if the hand has all specified cards, in the exact place
func (hand Hand) HasAllInPlace(kCards []Card) bool {
	if len(kCards) > len(hand) {
		return false
	}

	for i, k := range kCards {
		if hand[i] != k {
			return false
		}
	}
	return true
}

// SortForEnvido cmp func to sort Cards in a Hand, highest envido value first:
// 7-1,10,11,12
func SortForEnvido(a, b Card) int {
	an := a.N
	if an >= 10 {
		an = 0
	}
	bn := b.N
	if bn >= 10 {
		bn = 0
	}

	return int(bn) - int(an)
}

// SortForTruco cmp func to sort Cards in a Hand, highest truco value first.
// Uses TRUCO map in data.go
func SortForTruco(a, b Card) int {
	return int(b.Truco()) - int(a.Truco())
}

// Returns a sub-hand of the given hand
// of the cards that count for envido (2 or 1 card)
func (h Hand) EnvidoCards() *Hand {
	slices.SortFunc(h, SortForEnvido)

	s0 := h[0].S
	s1 := h[1].S
	s2 := h[2].S

	if s0 == s1 && s1 == s2 {
		// flor, highest envido for now
		return &Hand{h[0], h[1]}
	} else if s0 == s1 {
		return &Hand{h[0], h[1]}
	} else if s0 == s2 {
		return &Hand{h[0], h[2]}
	} else if s1 == s2 {
		return &Hand{h[1], h[2]}
	} else {
		return &Hand{h[0]}
	}
}

// Full value of hand
func (h Hand) Envido() uint8 {
	cards := h.EnvidoCards()

	if len(*cards) == 1 {
		return (*cards)[0].Envido()
	} else {
		return (*cards)[0].Envido() + (*cards)[1].Envido() + 20
	}
}

func (h Hand) Print() {
	for i := range len(h) {
		h[i].Print()
		fmt.Print(" ")
	}
}

func (h Hand) Println() {
	h.Print()
	fmt.Println()
}
