package main

func fact(x, l int) int {
	if x <= l {
		return l
	} else {
		return x * fact(x-1, l)
	}
}

// from n pick k
func pick(n, k int) float32 {
	return float32(fact(n, n-k+1)) / float32(fact(k, 1))
}

// Return all cards in `cards` that dont belong to list eCards
func CardsExcluding(cards, eCards []Card) []Card {
	cards_ := make([]Card, 0, len(cards))
	excluded := make(map[Card]bool)
	for _, c := range eCards {
		excluded[c] = true
	}
	for _, c := range cards {
		if !excluded[c] {
			cards_ = append(cards_, c)
		}
	}
	return cards_
}
