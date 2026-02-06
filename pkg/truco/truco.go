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

// IsReasonablyPlayed checks if the match between mHand and oHand is reasonably played.
//
// Strategies checked (for the player who plays second):
//   - Beat with minimum possible card that still wins
//   - Lose with minimum possible card when losing
//   - After a tie in round 0, play strongest card immediately
//   - Only tie in later rounds if already winning or cannot do better
//   - Must win R1 (with minimum possible) if already lost R0
//
// Returns false if either player's play violates these strategies (unreasonable/wasteful play).
func IsReasonablyPlayed(mHand, oHand Hand) bool {
	o0, o1, o2 := oHand[0].Truco(), oHand[1].Truco(), oHand[2].Truco()
	m0, m1, m2 := mHand[0].Truco(), mHand[1].Truco(), mHand[2].Truco()

	// Round 0: mHand plays first, oHand plays second (can strategize)
	if o0 > m0 {
		// oHand won - check minimum winning card
		if (o1 > m0 && o1 < o0) || (o2 > m0 && o2 < o0) {
			return false
		}
	} else if o0 < m0 {
		// oHand lost - check minimum losing card
		if (o1 < m0 && o1 < o0) || (o2 < m0 && o2 < o0) {
			return false
		}
	}
	// Tie in round 0 is acceptable

	// Round 1: winner of R0 plays first
	if o0 > m0 {
		// oHand won R0, plays first in R1, mHand can strategize
		// mHand lost R0, must try to win R1 if possible to stay alive
		if m1 > o1 {
			// mHand won R1 - check minimum winning card
			if m2 > o1 && m2 < m1 {
				return false
			}
		} else {
			// mHand lost or tied R1 - should have won if possible
			if m2 > o1 {
				return false
			}
			// Check minimum losing card -- to reduce double-counting
			if m2 < o1 && m2 < m1 {
				return false
			}
		}
	} else if o0 == m0 {
		// Tie in R0, mHand plays first in R1, oHand can strategize
		// After tie, should play strongest to maximize win chance
		if o2 > o1 {
			return false
		}
	} else {
		// mHand won R0, mHand plays first in R1, oHand can strategize
		// oHand must try to win R1 if possible to stay alive
		if o1 > m1 {
			// Winning R1 - check minimum winning card
			if o2 > m1 && o2 < o1 {
				return false
			}
		} else {
			// Losing or tying R1 after losing R0 (both bad)
			// Should have won if possible
			if o2 > m1 {
				return false
			}
			// If losing, check minimum loss -- to reduce double-counting
			if o1 < m1 && o2 < m1 && o2 < o1 {
				return false
			}
		}
	}

	// Round 2: only one card left, no choice to validate
	return true
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
	return float32(score) / math.PickC(37, 3)
}

type TrucoStats struct {
	StrengthAll  float32   // overall hand strength
	Count        int       // amount of hands simulated
	Perms        []Hand    // permutations of mHand
	StrengthPerm []float32 // strength, by permutations of mHand
	CountPerm    []float32 // amount of hands simulated, by permutations of mHand
}

func (stats TrucoStats) PPrint() {
	if stats.Count == 0 {
		if len(stats.Perms) == 0 {
			fmt.Println("Empty stats")
		} else {
			for _, hand := range stats.Perms {
				hand.Println()
			}
			fmt.Println("0 permutations simulated")
			fmt.Println("Hand is probably played sub-optimally")
		}
	} else {
		fmt.Println("Overall Strenght=", stats.StrengthAll)
		for i, hand := range stats.Perms {
			hand.Print()
			fmt.Printf(": strength=%.3f, of=%.0f\n", stats.StrengthPerm[i], stats.CountPerm[i])
		}
		fmt.Println(stats.Count, "permutations simulated")
	}
}

// TrucoStrengthStats calculates detailed strength statistics for a hand by simulating
// all possible permutations against all possible opponent hands, given known cards.
//
// Parameters:
//   - kCards: Cards held by the opponent (already played by them, in the order played).
//   - oCards: Known cards the opponent does not hold (e.g., cards played by other players).
//   - envido: Known envido of the opponent, helps exclude impossible hands.
//   - isMHandFirst: boolean controling who plays first in round 0
//   - hasStrategy: if true, checks IsReasonablyPlayed, discards permutations that are unreasonably played
//
// Notes:
//   - envido range as fsm envido (0-33: known, +100: range, 255: unknown)
//
// Returns TrucoStats containing the overall strength and per-permutation breakdown.
func (mHand Hand) TrucoStrengthStats(kCards, oCards []Card, envido uint8, isMHandFirst, hasStrategy bool) TrucoStats {
	mPerms := math.Permutations(mHand, 3)
	aCards := CardsExcluding(ALL_CARDS, append(mHand, oCards...))
	oPerms := math.Permutations(aCards, 3)

	isReasonablyPlayed := true
	var cScore, totScore, cCount, totCount int
	perms := make([]Hand, 0, 6)
	strengths := make([]float32, 0, 6)
	counts := make([]float32, 0, 6)

	for mH := range mPerms {
		cScore, cCount = 0, 0
		for oH := range oPerms {
			if envido != 255 {
				oEnvido := Hand(oH).Envido()
				if envido > 99 && oEnvido > (envido-100) { // range
					continue
				} else if envido < 99 && oEnvido != envido { // concrete
					continue
				}
			}

			if hasStrategy {
				if isMHandFirst {
					isReasonablyPlayed = IsReasonablyPlayed(mH, oH)
				} else {
					isReasonablyPlayed = IsReasonablyPlayed(oH, mH)
				}
			}

			if Hand(oH).HasAllInPlace(kCards) && isReasonablyPlayed {
				cScore += TrucoBeats(Hand(mH), Hand(oH))
				cCount++
			}
		}
		perms = append(perms, mH)
		strengths = append(strengths, (float32(cScore)/float32(cCount)+1)/2)
		counts = append(counts, float32(cCount))
		totScore += cScore
		totCount += cCount
	}

	return TrucoStats{
		StrengthAll:  (float32(totScore)/float32(totCount) + 1) / 2,
		Count:        totCount,
		Perms:        perms,
		StrengthPerm: strengths,
		CountPerm:    counts,
	}
}
