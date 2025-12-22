package main

import "fmt"

// We use uint8 for smallest size
// knowing that Cards will almost always
// align (in a Hand) in a single 8-byte block
type Card struct {
	n uint8 // Number: 1,2,3,4,5,6,7,10,11,12
	s uint8 // Suit: e,b,o,c (espada, basto, oro, copa)
}

// Value of card for envido (-20)
func (c Card) Envido() uint8 {
	if c.n <= 7 {
		return c.n
	} else {
		return 0
	}
}

func (c Card) Truco() uint8 {
	return TRUCO[c]
}

// Is a figure (10, 11, 12)
func (c Card) IsF() bool {
	return c.n >= 10
}

func (c Card) Print() {
	fmt.Printf("%d%c", c.n, c.s)
}

func (c Card) Println() {
	fmt.Printf("%d%c\n", c.n, c.s)
}
