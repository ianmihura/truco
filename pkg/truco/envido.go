package truco

import (
	"truco/pkg/math"
)

// returns a list of hands given a list of cards
// each hand consists of one card
func getListofOneCardHands(cards []Card) []Hand {
	hands := make([]Hand, len(cards))
	for i, c := range cards {
		hands[i] = Hand{c}
	}
	return hands
}

// list of pair (or single) cards that produce an envido score, eg.
// score=33, result={6e 7e} {6b 7b} {6o 7o} {6c 7c}
func envidoPairs(score uint8) (comb []Hand) {
	// assert score >= 0 and score <= 33, "Envido must be between 0 and 33"
	cardVal := ENVIDOS[score]

	if score < 20 {
		// single card
		if score == 0 {
			comb = getListofOneCardHands(FIGURES)
		} else {
			comb = getListofOneCardHands(GetSuitCards(cardVal[0][0]))
		}
	} else if score == 20 {
		// two figures
		for cards := range math.Combinations(FIGURES, 2) {
			if cards[0].S == cards[1].S {
				comb = append(comb, Hand(cards))
			}
		}
	} else {
		// normal envido
		for v := range cardVal {
			var cardsToComb []Card
			if cardVal[v][0] == 'f' {
				// has a figure
				cardsToComb = append(FIGURES, GetSuitCards(cardVal[v][1])...)
			} else {
				cardsToComb = append(GetSuitCards(cardVal[v][0]), GetSuitCards(cardVal[v][1])...)
			}

			for cards := range math.Combinations(cardsToComb, 2) {
				if cards[0].S == cards[1].S { // same suit
					if !cards[0].IsF() || !cards[1].IsF() { // exclude double figures
						comb = append(comb, Hand(cards))
					}
				}
			}
		}
	}
	return comb
}

// gets all possible hands given an envido
func EnvidoHands(score uint8) (hands []Hand) {
	pairs := envidoPairs(score)

	for _, p := range pairs {
		if len(p) == 2 {
			// p has 2 cards of same suit.
			// we need 1 more card of DIFFERENT suit.
			// rank does not matter for the 3rd card as it doesn't affect envido score
			// (envido score is determined by the 2 same-suit cards).
			suit := p[0].S
			for _, c := range ALL_CARDS {
				if c.S != suit {
					// Create new hand with the pair + this card
					h := make(Hand, 3)
					h[0] = p[0]
					h[1] = p[1]
					h[2] = c
					hands = append(hands, h)
				}
			}
		} else if len(p) == 1 {
			// p has 1 card.
			// we need 2 more cards.
			// To maintain the score, the other 2 cards must NOT form a pair with each other or with p[0]
			// that results in a higher score OR a score >= 20.
			// Basically, for single card score < 20, it means NO pairs are allowed in the hand.
			// So all 3 cards must be different suits.
			// Also, to ensure p[0] is the one providing the score, the other cards must have
			// envido value <= p[0].Envido().

			c1 := p[0]

			// Filter candidates: different suit than c1, and value <= c1.Envido()
			var candidates []Card
			for _, c := range ALL_CARDS {
				if c.S != c1.S && c.Envido() <= c1.Envido() {
					candidates = append(candidates, c)
				}
			}

			// Get combinations of 2 from candidates
			for combo := range math.Combinations(candidates, 2) {
				// We need to ensure the 2 candidates are ALSO different suits from each other
				if combo[0].S != combo[1].S {
					h := make(Hand, 3)
					h[0] = c1
					h[1] = combo[0]
					h[2] = combo[1]
					hands = append(hands, h)
				}
			}
		}
	}
	return hands
}

// Raw probability of having a certain envido score
func PEnvido(score uint8) float32 {
	return 0.0 // TODO
}

// Probability a given envido is the highest of the table, given mCards and other kCards
func PEnvidoHighest(score uint8, mCards, kCards []Card) float32 {
	return 0.0 // TODO
}
