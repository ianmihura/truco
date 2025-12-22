package main

import (
	"slices"
)

// Simulates two hands being played in Truco, in the orther they were given
//
// returns:
//   - 1 if mHand beats oHand
//   - -1 if mHand looses against oHand
//   - 0 if there's a tie (usually the first player wins)
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

// returns true if my sorted hand beats the other sorted hand, or if there's
// a full tie (assumes i am 'mano': first player)
func (mHand Hand) TrucoBeatsSorted(oHand Hand) bool {
	slices.SortFunc(mHand, sortForTruco)
	slices.SortFunc(oHand, sortForTruco)

	return TrucoBeats(mHand, oHand) >= 0
}

// returns count of beats-losses of mHand against oHand
func (mHand Hand) TrucoBeatsAll(oHand Hand) (score int) {
	mPerms := Permutations(mHand, 3)
	oPerms := Permutations(oHand, 3)
	for mH := range mPerms {
		for oH := range oPerms {
			score += TrucoBeats(mH, oH)
		}
	}
	return score
}

// strength of a hand in truco
//
// plays the hand against all other hands, in all possible permutations.
// counts times it wins, minus losses. averages the result dividing by 36:
// range of score = (-36 to 36)
func (mHand Hand) TrucoStrength() float32 {
	mPerms := Permutations(mHand, 3)
	aCards := CardsExcluding(ALL_CARDS, mHand)
	oPerms := Permutations(aCards, 3)
	var score int
	for mH := range mPerms {
		for oH := range oPerms {
			score += TrucoBeats(mH, oH)
		}
	}
	return float32(score) / pick(37, 3)
}
