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
//   - 0 if there's a tie or loss
func TrucoBeats(mHand, oHand Hand, m Card) int {
	var s0, s1, s2 int
	var o0, o1, o2, m0, m1, m2 uint8

	if m == NO_CARD {
		o0, o1, o2 = oHand[0].Truco(), oHand[1].Truco(), oHand[2].Truco()
		m0, m1, m2 = mHand[0].Truco(), mHand[1].Truco(), mHand[2].Truco()
	} else {
		o0, o1, o2 = oHand[0].TrucoUY(m), oHand[1].TrucoUY(m), oHand[2].TrucoUY(m)
		m0, m1, m2 = mHand[0].TrucoUY(m), mHand[1].TrucoUY(m), mHand[2].TrucoUY(m)
	}

	if m0 > o0 {
		s0 = 1
	} else if m0 < o0 {
		s0 = -1
	}
	if m1 > o1 {
		s1 = 1
	} else if m1 < o1 {
		s1 = -1
	}
	if m2 > o2 {
		s2 = 1
	} else if m2 < o2 {
		s2 = -1
	}

	var res int
	if s0 == 0 {
		// tie in the first round is defined inmediately after
		if s1 == 0 {
			res = s2
		} else {
			res = s1
		}

	} else if s1 == 0 {
		// any other tie is defined by winner of first round
		res = s0

	} else if s0 == s1 {
		// a player won first two rounds
		res = s0

	} else if s2 == 0 {
		// after we are sure there's no winner of the first two rounds,
		// (first two round alternate winners)
		// tie in last round is defined by winner of first round
		res = s0

	} else {
		// first two round alternate winners, last round defines
		res = s2
	}

	if res == 1 {
		return 1
	} else {
		return 0
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
func IsReasonablyPlayed(mHand, oHand Hand, m Card) bool {
	var o0, o1, o2, m0, m1, m2 uint8

	if m == NO_CARD {
		o0, o1, o2 = oHand[0].Truco(), oHand[1].Truco(), oHand[2].Truco()
		m0, m1, m2 = mHand[0].Truco(), mHand[1].Truco(), mHand[2].Truco()
	} else {
		o0, o1, o2 = oHand[0].TrucoUY(m), oHand[1].TrucoUY(m), oHand[2].TrucoUY(m)
		m0, m1, m2 = mHand[0].TrucoUY(m), mHand[1].TrucoUY(m), mHand[2].TrucoUY(m)
	}

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

// strength of a hand in truco (brute force)
//
// plays the hand against all other hands, in all possible permutations.
// counts times it wins, minus losses. Normalizes result to a percent.
// range of score = (0 to 1)
func (mHand Hand) TrucoStrength() float32 {
	mPerms := math.PermutationsRaw(mHand, 3)
	aCards := CardsExcluding(ALL_CARDS, mHand)
	oPerms := math.PermutationsRaw(aCards, 3)
	var score int
	for mH := range mPerms {
		for oH := range oPerms {
			score += TrucoBeats(Hand(mPerms[mH]), Hand(oPerms[oH]), NO_CARD)
		}
	}
	return float32(score) / (math.PickC(37, 3) * 36.0)
}

// strength of a hand in truco uruguay (brute force)
//
// plays the hand against all other hands, in all possible permutations, with all possible muestras.
// counts times it wins, minus losses. Normalizes result to a percent.
// range of score = (0 to 1)
//
// bench = 330 ms
func (mHand Hand) TrucoStrengthUY() float32 {
	var c int
	mPerms := math.PermutationsRaw(mHand, 3)
	aCards := CardsExcluding(ALL_CARDS, mHand)
	oPerms := math.PermutationsRaw(aCards, 3)
	pPerms := math.PermutationsRaw(aCards, 1) // muestra

	var score int
	for mH := range mPerms {
		for oH := range oPerms {
			for m := range pPerms {
				if oPerms[oH][0] == pPerms[m][0] || oPerms[oH][1] == pPerms[m][0] || oPerms[oH][2] == pPerms[m][0] {
					// muestra may be in oHand
					continue // muestra should be unique
				} else {
					score += TrucoBeats(Hand(mPerms[mH]), Hand(oPerms[oH]), pPerms[m][0])
					c++
				}
			}
		}
	}
	return float32(score) / float32(c)
}

// TrucoStrengthStats calculates strength statistics for a mHand by simulating
// all possible permutations against all possible opponent hands, given known info.
// Helps players identify best permutation to play
//
// For Argentinian Truco.
//
// Parameters:
//   - kCards: Cards held by the opponent (already played by them, in the order played).
//   - oCards: Known cards the opponent does not hold (e.g., cards played by other players).
//   - envido: Known envido of the opponent, helps exclude impossible hands.
//   - isMHandFirst: boolean controling who plays first in round 0
//   - hasStrategy: if true, checks IsReasonablyPlayed, discards permutations that are unreasonably played
//
// Notes:
//   - envido range as fsm envido (0-33: known, 100-199: envido range, 200: unknown flor, 201-254: flor, 255: unknown)
//
// Returns TrucoStats containing the overall strength and per-permutation breakdown.
func (mHand Hand) TrucoStrengthStats(kCards, oCards []Card, envido uint8, isMHandFirst, hasStrategy bool) TrucoStats {
	mPerms := math.Permutations(mHand, 3)
	aCards := CardsExcluding(ALL_CARDS, append(mHand, oCards...))
	oPerms := math.Permutations(aCards, 3)
	mEnvido := mHand.Envido()

	isReasonablyPlayed := true
	var eScore, eCount, cScore, totScore, cCount, totCount int
	perms := make([]Hand, 0, 6)
	winsPerm := make([]float32, 0, 6)
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

			if Hand(oH).HasAllInPlace(kCards) {
				if hasStrategy {
					if isMHandFirst {
						isReasonablyPlayed = IsReasonablyPlayed(mH, oH, NO_CARD)
					} else {
						isReasonablyPlayed = IsReasonablyPlayed(oH, mH, NO_CARD)
					}
				}

				if isReasonablyPlayed {
					cScore += TrucoBeats(Hand(mH), Hand(oH), NO_CARD)
					cCount++
				}
				eScore += EnvidoBeats(mEnvido, Hand(oH).Envido(), isMHandFirst)
				eCount++
			}
		}
		perms = append(perms, mH)
		winsPerm = append(winsPerm, float32(cScore))
		counts = append(counts, float32(cCount))
		totScore += cScore
		totCount += cCount
	}

	return finalTrucoStrengthStats(rawTrucoStats{
		TotCount: totCount,
		TotScore: totScore,
		WinsPerm: winsPerm,
		Counts:   counts,
		MHand:    mHand,
		Perms:    perms,
		MEnvido:  mEnvido,
		EScore:   eScore,
		ECount:   eCount,
	})
}

// TrucoStrengthStatsUY calculates strength statistics for a mHand by simulating
// all possible permutations against all possible opponent hands, given known info.
// Helps players identify best permutation to play.
//
// For Uruguayan Truco.
//
// Parameters:
//   - kCards: Cards held by the opponent (already played by them, in the order played).
//   - oCards: Known cards the opponent does not hold (e.g., cards played by other players). First card is 'muestra'.
//   - envido: Known envido of the opponent, helps exclude impossible hands.
//   - isMHandFirst: boolean controling who plays first in round 0
//   - hasStrategy: if true, checks IsReasonablyPlayed, discards permutations that are unreasonably played
//
// Notes:
//   - envido range as fsm envido (0-33: known, 100-199: envido range, 200: unknown flor, 201-254: flor, 255: unknown)
//
// Returns TrucoStats containing the overall strength and per-permutation breakdown.
func (mHand Hand) TrucoStrengthStatsUY(kCards, oCards []Card, envido uint8, isMHandFirst, hasStrategy bool) TrucoStats {
	mPerms := math.PermutationsRaw(mHand, 3)
	aCards := CardsExcluding(ALL_CARDS, append(mHand, oCards...))
	oPerms := math.PermutationsRaw(aCards, 3)
	muestra := oCards[0]
	mEnvido := mHand.EnvidoUY(muestra)

	isReasonablyPlayed := true
	var eScore, eCount, cScore, totScore, cCount, totCount int
	perms := make([]Hand, 0, 6)
	winsPerm := make([]float32, 0, 6)
	counts := make([]float32, 0, 6)

	for mH := range mPerms {
		cScore, cCount = 0, 0
		for oH := range oPerms {
			if oPerms[oH][0] == muestra || oPerms[oH][1] == muestra || oPerms[oH][2] == muestra {
				continue // muestra should be unique
			}

			oEnvido := Hand(oPerms[oH]).EnvidoUY(muestra)
			if envido == 255 { // didnt declare anything
				if oEnvido > 200 { // only filter out flor
					continue
				}
			} else if envido == 200 { // declare unknown flor
				if oEnvido < 200 || oEnvido == 255 { // filter out non-flor
					continue
				}
			} else if envido < 99 { // declare concrete envido
				if oEnvido != envido {
					continue
				}
			} else if envido < 199 { // declare range envido 'son buenas'
				if oEnvido > (envido - 100) { // range
					continue
				}
			} else { // known flor
				if oEnvido != envido {
					continue
				}
			}

			if Hand(oPerms[oH]).HasAllInPlace(kCards) {
				if hasStrategy {
					if isMHandFirst {
						isReasonablyPlayed = IsReasonablyPlayed(mPerms[mH], oPerms[oH], muestra)
					} else {
						isReasonablyPlayed = IsReasonablyPlayed(oPerms[oH], mPerms[mH], muestra)
					}
				}

				if isReasonablyPlayed {
					cScore += TrucoBeats(Hand(mPerms[mH]), Hand(oPerms[oH]), muestra)
					cCount++
				}
				eScore += EnvidoBeats(mEnvido, oEnvido, isMHandFirst)
				eCount++
			}
		}
		perms = append(perms, mPerms[mH])
		winsPerm = append(winsPerm, float32(cScore))
		counts = append(counts, float32(cCount))
		totScore += cScore
		totCount += cCount
	}

	return finalTrucoStrengthStats(rawTrucoStats{
		TotCount: totCount,
		TotScore: totScore,
		WinsPerm: winsPerm,
		Counts:   counts,
		MHand:    mHand,
		Perms:    perms,
		MEnvido:  mEnvido,
		EScore:   eScore,
		ECount:   eCount,
	})
}

type rawTrucoStats struct {
	TotCount int
	TotScore int
	WinsPerm []float32
	Counts   []float32
	MHand    Hand
	Perms    []Hand
	MEnvido  uint8
	EScore   int
	ECount   int
}

// finalTrucoStrengthStats calculates the final Stats from the raw simulation results.
func finalTrucoStrengthStats(rawStats rawTrucoStats) TrucoStats {
	var strengthAll float32
	var strengthsPerm []float32
	var strengthsPermRel []float32
	var strengthsPosition [][]float32

	if rawStats.TotCount > 0 {
		strengthAll = float32(rawStats.TotScore) / float32(rawStats.TotCount)
		for i := range rawStats.WinsPerm {
			strengthsPerm = append(strengthsPerm, rawStats.WinsPerm[i]/float32(rawStats.TotCount))
			strengthsPermRel = append(strengthsPermRel, rawStats.WinsPerm[i]/rawStats.Counts[i])
		}
	} else {
		strengthAll = 0
		strengthsPerm = []float32{0, 0, 0, 0, 0, 0}
		strengthsPermRel = []float32{0, 0, 0, 0, 0, 0}
	}

	for pos := range 3 {
		posScoreArr := make([]float32, 0, 3)
		for _, mCard := range rawStats.MHand {
			var pScore float32
			for iPerm, perm := range rawStats.Perms {
				if mCard == perm[pos] {
					pScore += strengthsPerm[iPerm]
				}
			}
			posScoreArr = append(posScoreArr, pScore)
		}
		strengthsPosition = append(strengthsPosition, posScoreArr)
	}

	sMHand := make([]string, 0, len(rawStats.MHand))
	for c := range rawStats.MHand {
		sMHand = append(sMHand, rawStats.MHand[c].ToEmoji())
	}

	var mEnvidoScore float32
	if rawStats.ECount > 0 {
		mEnvidoScore = float32(rawStats.EScore) / float32(rawStats.ECount)
	}

	return TrucoStats{
		MHand:            sMHand,
		StrengthAll:      strengthAll,
		Count:            rawStats.TotCount,
		Perms:            rawStats.Perms,
		WinsPerm:         rawStats.WinsPerm,
		StrengthPermRel:  strengthsPermRel,
		StrengthPermAbs:  strengthsPerm,
		StrengthPosition: strengthsPosition,
		CountPerm:        rawStats.Counts,
		MEnvido:          rawStats.MEnvido,
		MEnvidoScore:     mEnvidoScore,
	}
}

type TrucoStats struct {
	MHand            []string    // my hand: input parameter
	StrengthAll      float32     // overall hand strength: % hands you win
	Count            int         // amount of hands simulated
	Perms            []Hand      // permutations of mHand
	WinsPerm         []float32   // hands you win
	StrengthPermRel  []float32   // strength: hands you win / hands played with this perm
	StrengthPermAbs  []float32   // strength: hands you win / hands played in total
	StrengthPosition [][]float32 //
	CountPerm        []float32   // amount of hands simulated, by permutations of mHand
	MEnvido          uint8       // my envido
	MEnvidoScore     float32     // my envido strength: hands you win / hands played in total
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
			fmt.Printf(": strength=%.3f, of=%.0f\n", stats.WinsPerm[i], stats.CountPerm[i])
		}
		fmt.Println(stats.Count, "permutations simulated")
	}
}
