package truco

import (
	"fmt"
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

type TrucoStats struct {
	StrengthAll  float32
	Perms        []Hand
	StrengthPerm []float32
	Count        int
}

func (stats TrucoStats) PPrint() {
	fmt.Println(stats.StrengthAll)
	for i, hand := range stats.Perms {
		hand.Print()
		fmt.Println(": strength", stats.StrengthPerm[i])
	}
	fmt.Println("Played", stats.Count, "permutations")
}

// TrucoStrengthStats calculates detailed strength statistics for a hand by simulating
// all possible permutations against all possible opponent hands, given known cards.
//
// Parameters:
//   - oCards: Cards held by the opponent (already played by them).
//   - kCards: Known cards to exclude from the deck (e.g., cards played by other players).
//
// Returns TrucoStats containing the overall strength and per-permutation breakdown.
func (mHand Hand) TrucoStrengthStats(oCards, kCards []Card) TrucoStats {
	mPerms := math.Permutations(mHand, 3)
	aCards := CardsExcluding(ALL_CARDS, append(mHand, kCards...))
	oPerms := math.Permutations(aCards, 3)

	var cScore, totScore, count int
	strengths := make([]float32, 0, 6)
	perms := make([]Hand, 0, 6)

	for mH := range mPerms {
		cScore = 0
		perms = append(perms, mH)
		for oH := range oPerms {
			if Hand(oH).HasAll(oCards) {
				cScore += TrucoBeats(Hand(mH), Hand(oH))
				count++
			}
		}
		strengths = append(strengths, float32(cScore))
		totScore += cScore
	}

	for i, s := range strengths {
		strengths[i] = (s/(float32(count)/6) + 1) / 2
	}

	return TrucoStats{
		StrengthAll:  (float32(totScore)/float32(count) + 1) / 2,
		Perms:        perms,
		StrengthPerm: strengths,
		Count:        count,
	}
}
