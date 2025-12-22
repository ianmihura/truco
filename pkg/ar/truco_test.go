package ar

import "testing"

func TestTrucoBeats(t *testing.T) {
	// Cases from the request
	// 1 1 x => 1
	checkTrucoBeats(t, 1, 1, 2, 1) // 2 is dummy for x: any other value
	checkTrucoBeats(t, 1, 0, 2, 1)
	checkTrucoBeats(t, 1, -1, 1, 1)
	checkTrucoBeats(t, 1, -1, -1, -1)
	checkTrucoBeats(t, 1, -1, 0, 1)

	// -1 -1 x => -1
	checkTrucoBeats(t, -1, -1, 2, -1)
	checkTrucoBeats(t, -1, 0, 2, -1)
	checkTrucoBeats(t, -1, 1, -1, -1)
	checkTrucoBeats(t, -1, 1, 1, 1)
	checkTrucoBeats(t, -1, 1, 0, -1)

	// 0 1 x => 1
	checkTrucoBeats(t, 0, 1, 2, 1)
	checkTrucoBeats(t, 0, -1, 2, -1)
	checkTrucoBeats(t, 0, 0, 1, 1)
	checkTrucoBeats(t, 0, 0, -1, -1)
	checkTrucoBeats(t, 0, 0, 0, 0)
}

func checkTrucoBeats(t *testing.T, s0, s1, s2, expected int) {
	// Helper to construct hands that produce the desired score array
	// Truco values logic:
	// mHand[i] > oHand[i] -> 1
	// mHand[i] < oHand[i] -> -1
	// mHand[i] == oHand[i] -> 0

	mHand := make(Hand, 3)
	oHand := make(Hand, 3)

	setupRound(s0, 0, &mHand, &oHand)
	setupRound(s1, 1, &mHand, &oHand)
	setupRound(s2, 2, &mHand, &oHand)

	result := TrucoBeats(mHand, oHand)
	if result != expected {
		t.Errorf("For scores [%d, %d, %d], expected %d, got %d", s0, s1, s2, expected, result)
	}
}

func setupRound(s, i int, mHand, oHand *Hand) {
	// Use arbitrary card values
	high := Card{1, 'e'} // Value 14
	low := Card{4, 'c'}  // Value 1
	mid := Card{1, 'c'}  // Value 8

	// Adjust values to ensure strict > or < if not 0
	// Actually we can just trick the TrucoMap or use known values.

	switch s {
	case 1:
		(*mHand)[i] = high
		(*oHand)[i] = low
	case -1:
		(*mHand)[i] = low
		(*oHand)[i] = high
	default:
		(*mHand)[i] = mid
		(*oHand)[i] = mid
	}
}

func TestTrucoBeatsAll(t *testing.T) {
	mHand := Hand{{1, 'e'}, {1, 'b'}, {7, 'e'}}
	oHand := Hand{{10, 'e'}, {10, 'c'}, {4, 'e'}}
	score := mHand.TrucoBeatsAll(oHand)
	if score != 36 {
		t.Errorf("{1, 'e'}, {1, 'b'}, {7, 'e'} should always beat any other hand, instead got %d", score)
	}

	mHand = Hand{{10, 'b'}, {10, 'o'}, {4, 'b'}}
	score = mHand.TrucoBeatsAll(oHand)
	if score != 0 {
		t.Errorf("equal hands should have score=0, instead got %d", score)
	}

	// only one card is better than the other hand
	// for every order i can play, there's only 2 orders that i beat the oponent:
	//  - we must play our different card in the same turn (eg. first card played is 4 and 5 respectively)
	//    + 3 positions x2 permutations of my other two cards
	//  - the oponent can permute the other two cards
	//    + 3 positions x2 permutations
	// total = 3x2 + 3x2 = 12
	mHand = Hand{{10, 'b'}, {10, 'o'}, {5, 'b'}}
	score = mHand.TrucoBeatsAll(oHand)
	if score != 12 {
		t.Errorf("only one card different should beat 12 hands, instead got %d", score)
	}
}
