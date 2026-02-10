package truco

import (
	"fmt"
)

/*
 */

// We use uint8 for smallest size
// knowing that Cards will almost always
// align (in a Hand) in a single 8-byte block
//
// Types of Cards/Hands:
//   - Flat Cards/Hands are hands in truco argentino, and hands in uruguay without muestra.
//   - Uruguay Cards/Hands may include 'p' as a suit (pieza).
type Card struct {
	N uint8 // Number: 1,2,3,4,5,6,7,10,11,12.
	S uint8 // Suit: e,b,o,c (espada, basto, oro, copa). For uruguay: p=pieza.
}

// Returns a new Card from a string: Card{c[0], c[1]}.
//
//	"1e" -> Card{1, 'e'}
//	"10e" -> Card{10, 'e'}
func NewCard(c string) Card {
	if len(c) == 2 {
		return Card{c[0] - '0', c[1]}
	} else if len(c) == 3 {
		return Card{10 + (c[1] - '0'), c[2]}
	} else {
		return Card{}
	}
}

// Converts in place a flat card to a uruguay card, given a m=muestra Card.
// Is destructive: will overwrite the value of the suit for 'p' where necesary.
func (c *Card) UY(m Card) {
	// we do it this extended way for efficiency
	if c.S == m.S {
		if 2 == c.N || 4 == c.N || 5 == c.N || 10 == c.N || 11 == c.N {
			c.S = 'p'
		} else if c.N == 12 && (2 == m.N || 4 == m.N || 5 == m.N || 10 == m.N || 11 == m.N) {
			c.N = m.N
			c.S = 'p'
		}
	}
}

// Value of card for envido (-20)
func (c Card) Envido() uint8 {
	if c.N <= 7 {
		return c.N
	} else {
		return 0
	}
}

func (c Card) Truco() uint8 {
	return GetTruco(c)
}

func (c Card) TrucoUY(m Card) uint8 {
	c.UY(m)
	return GetTruco(c)
}

// Is a figure (10, 11, 12)
func (c Card) IsF() bool {
	return c.N >= 10
}

// string representation of card
func (c Card) ToString() string {
	return fmt.Sprintf("%d%c", c.N, c.S)
}

// Returns the rank of the card. Eg.
//   - 1e -> 1e
//   - 2c -> 2
//   - 7c -> 7f
func (c Card) ToRank() string {
	if c.Truco() > 10 {
		// 1e, 1b, 7e, 7b
		return fmt.Sprintf("%d%c", c.N, c.S)
	} else if c.N == 1 || c.N == 7 {
		// 1f, 7f
		return fmt.Sprintf("%df", c.N)
	} else {
		return fmt.Sprintf("%d", c.N)
	}
}

func (c Card) Print() {
	fmt.Printf("%d%c", c.N, c.S)
}

func (c Card) Println() {
	fmt.Printf("%d%c\n", c.N, c.S)
}

// Given a list of cards, returns a list of non-null cards
func RealCards(cards []Card) []Card {
	cards_ := make([]Card, 0, len(cards))
	for c := range cards {
		if cards[c].N != 0 {
			cards_ = append(cards_, cards[c])
		}
	}
	return cards_
}

// returns a list of all suits of the same n, eg.
// n=1, returns []Card{1e, 1b, 1o, 1c}
func GetSuitCards(n uint8) []Card {
	return []Card{{n, 'e'}, {n, 'b'}, {n, 'o'}, {n, 'c'}}
}
