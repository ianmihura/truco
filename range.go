package main

import "slices"

// Returns true only if all `cards` belong in `iCards`
func isEveryCardIncluded(cards, iCards []Card) bool {
	for _, card := range cards {
		if !slices.Contains(iCards, card) {
			return false
		}
	}
	return true
}

// Given cards player holds (mCards) and list of all possible cards they could hold (aCards),
// returns a list of possible hands they could have
//
// Consider that:
// - if len(mCards) == 3, then len(hands) == 1
// - if len(mCards) == 0, then len(hands) == pick(aCards, 3)
func cardRangeNoEnvido(aCards, mCards []Card) []Hand {
	hands := make([]Hand, 0, int(pick(len(aCards), 3)))

	var hand Hand
	combo := Combinations(aCards, 3)
	for cs := range combo {
		if isEveryCardIncluded(mCards, cs) {
			hand = Hand{}
			for _, c := range cs {
				hand = append(hand, c)
			}
			hands = append(hands, hand)
		}
	}
	return hands
}

// Given an envido score, cards player holds (mCards) and other known cards they dont (kCards),
// whats a list of possible hands they could have
//
// Consider that:
// - if len(mCards) == 3, then len(hands) == 1
// - as len(kCards) grows, len(hands) shrinks
// - len(hands) is not homogeneous over all envido scores
func CardRange(score uint8, mCards, kCards []Card) []Hand {
	aCards := CardsExcluding(ALL_CARDS, kCards)
	hands_ := cardRangeNoEnvido(aCards, mCards)
	hands := make([]Hand, 0, len(hands_))

	if score == 255 {
		// unknown envido
		return hands_
	} else if score >= 100 {
		// only know that envido <= (score-100)
		score_ := score - 100
		for _, h := range hands_ {
			if h.Envido() <= score_ {
				hands = append(hands, h)
			}
		}
		return hands
	} else {
		for _, h := range hands_ {
			if h.Envido() == score {
				hands = append(hands, h)
			}
		}
	}
	return hands
}
