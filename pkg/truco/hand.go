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
	if handStr == "" {
		return []Card{}
	}
	cards := strings.Split(handStr, " ")
	hand := make([]Card, 0, len(cards))
	for c := range cards {
		hand = append(hand, NewCard(cards[c]))
	}
	return Hand(hand)
}

// Converts a flat hand to a uruguay hand, given a m=muestra Card.
// Returns a new hand, does not overwrite the old one.
func (hand Hand) UY(m Card) Hand {
	newHand := make(Hand, 3)
	copy(newHand, hand)
	newHand[0].UY(m)
	newHand[1].UY(m)
	newHand[2].UY(m)
	return newHand
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

// Full value of hand, 0-33
func (h Hand) Envido() uint8 {
	slices.SortFunc(h, SortForEnvido)

	s0, s1, s2 := h[0].S, h[1].S, h[2].S

	if s0 == s1 && s1 == s2 {
		// flor, highest envido for now
		return h[0].Envido() + h[1].Envido() + 20
	} else if s0 == s1 {
		return h[0].Envido() + h[1].Envido() + 20
	} else if s0 == s2 {
		return h[0].Envido() + h[2].Envido() + 20
	} else if s1 == s2 {
		return h[1].Envido() + h[2].Envido() + 20
	} else {
		return h[0].Envido()
	}
}

// Full value of hand, including flor.
//
//   - 0-37 envido
//   - 220-247 flor
func (hand Hand) EnvidoUY(m Card) uint8 {
	var pCards []Card
	var normalCards []Card

	h := hand.UY(m)

	// count piezas
	for i := range 3 {
		c := h[i]
		if c.S == 'p' {
			pCards = append(pCards, c)
		} else {
			normalCards = append(normalCards, c)
		}
	}

	if len(pCards) == 0 {
		// No piezas, get envido score
		slices.SortFunc(h, SortForEnvido)

		s0, s1, s2 := h[0].S, h[1].S, h[2].S

		if s0 == s1 && s1 == s2 {
			// flor all same suit
			return h[0].Envido() + h[1].Envido() + h[2].Envido() + 220
		} else if s0 == s1 {
			return h[0].Envido() + h[1].Envido() + 20
		} else if s0 == s2 {
			return h[0].Envido() + h[2].Envido() + 20
		} else if s1 == s2 {
			return h[1].Envido() + h[2].Envido() + 20
		} else {
			return h[0].Envido()
		}
	}

	if len(pCards) == 1 {
		if normalCards[0].S == normalCards[1].S {
			// flor: 1 pieza + envido
			return ENVIDO_PIEZA[pCards[0].N] + normalCards[0].Envido() + normalCards[1].Envido() + 220

		} else {
			// envido: 1 pieza + carta alta
			slices.SortFunc(normalCards, SortForEnvido)
			return ENVIDO_PIEZA[pCards[0].N] + normalCards[0].Envido() + 20
		}
	}

	if len(pCards) == 2 {
		// flor: 2 piezas
		sum := uint8(0)
		sum += ENVIDO_PIEZA[pCards[0].N]
		sum += ENVIDO_PIEZA[pCards[1].N]
		sum += normalCards[0].Envido()
		return sum + 220
	}

	if len(pCards) == 3 {
		// flor: 3 piezas
		sum := uint8(0)
		sum += ENVIDO_PIEZA[pCards[0].N]
		sum += ENVIDO_PIEZA[pCards[1].N]
		sum += ENVIDO_PIEZA[pCards[2].N]
		return sum + 220
	}

	return 255 // should never get here
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

func (h Hand) ToString() string {
	var result strings.Builder
	for c := range h {
		result.WriteString(" " + h[c].ToString())
	}
	return result.String()
}

func (h Hand) ToEmoji() string {
	var result strings.Builder
	for c := range h {
		result.WriteString(" " + h[c].ToEmoji())
	}
	return result.String()
}
