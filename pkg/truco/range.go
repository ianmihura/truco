package truco

import (
	"slices"
	"truco/pkg/math"
)

func CardsLowerEqual(min uint8) []string {
	result := make([]string, 0, len(RANKS))
	for rank := range RANKS {
		if RANKS[rank] <= min {
			result = append(result, rank)
		}
	}
	return result
}

// Return all cards in `cards` that dont belong to list eCards
func CardsExcluding(cards, eCards []Card) []Card {
	if len(eCards) == 0 {
		return cards
	}

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
	hands := make([]Hand, 0, int(math.PickC(len(aCards), 3)))

	combo := math.Combinations(aCards, 3)
	for cs := range combo {
		if isEveryCardIncluded(mCards, cs) {
			hands = append(hands, Hand(cs))
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
