package truco

import (
	"truco/pkg/math"
)

// Simulates two hands being played in Truco,
// cards are played in the orther they were given (by index)
//
// returns:
//   - 1 if mHand beats oHand
//   - -1 if mHand looses against oHand
//   - 0 if there's a tie
func TrucoBeats(mHand, oHand Hand) int {
	score := make([]int, 3)
	for i := range 3 {
		if mHand[i].Truco() > oHand[i].Truco() {
			score[i] = 1
		} else if mHand[i].Truco() < oHand[i].Truco() {
			score[i] = -1
		} else {
			score[i] = 0
		}
	}

	s0, s1, s2 := score[0], score[1], score[2]
	if s0 == 0 {
		// tie in the first round is defined inmediately after
		if s1 == 0 {
			return s2
		} else {
			return s1
		}

	} else if s1 == 0 {
		// any other tie is defined by winner of first round
		return s0

	} else if s0 == s1 {
		// a player won first two rounds
		return s0

	} else if s2 == 0 {
		// after we are sure there's no winner of the first two rounds,
		// (first two round alternate winners)
		// tie in last round is defined by winner of first round
		return s0

	} else {
		// first two round alternate winners, last round defines
		return s2
	}
}

// returns count of beats-losses of all permutations of mHand against oHand
func (mHand Hand) TrucoBeatsAll(oHand Hand) (score int) {
	mPerms := math.Permutations(mHand, 3)
	oPerms := math.Permutations(oHand, 3)
	for mH := range mPerms {
		for oH := range oPerms {
			score += TrucoBeats(Hand(mH), Hand(oH))
		}
	}
	return score
}

// strength of a hand in truco (brute force)
//
// plays the hand against all other hands, in all possible permutations.
// counts times it wins, minus losses. averages the result dividing by 36:
// range of score = (-36 to 36)
func (mHand Hand) TrucoStrength() float32 {
	mPerms := math.Permutations(mHand, 3)
	aCards := CardsExcluding(ALL_CARDS, mHand)
	oPerms := math.Permutations(aCards, 3)
	var score int
	for mH := range mPerms {
		for oH := range oPerms {
			score += TrucoBeats(Hand(mH), Hand(oH))
		}
	}
	return float32(score) / math.Pick(37, 3)
}
