package ar

import (
	"fmt"
	"slices"
)

// Slice of Cards of any length
type Hand []Card

// cmp func to sort Cards in a Hand, highest envido value first:
// 7-1,10,11,12
func sortForEnvido(a, b Card) int {
	an := a.n
	if an >= 10 {
		an = 0
	}
	bn := b.n
	if bn >= 10 {
		bn = 0
	}

	return int(bn) - int(an)
}

// cmp func to sort Cards in a Hand, highest truco value first.
// Uses TRUCO map in data.go
func sortForTruco(a, b Card) int {
	return int(b.Truco()) - int(a.Truco())
}

// Returns a sub-hand of the given hand
// of the cards that count for envido (2 or 1 card)
func (h Hand) EnvidoCards() *Hand {
	slices.SortFunc(h, sortForEnvido)

	s0 := h[0].s
	s1 := h[1].s
	s2 := h[2].s

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

// Full value of hand
func (h Hand) Truco() (s uint8) {
	for _, c := range h {
		s += c.Truco()
	}
	return s
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
