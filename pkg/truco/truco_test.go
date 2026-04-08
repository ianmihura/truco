package truco

import (
	"testing"
)

func TestTrucoBeats(t *testing.T) {
	// Cases from the request
	// 1 1 x => 1
	checkTrucoBeats(t, 1, 1, 2, 1) // 2 is dummy for x: any other value
	checkTrucoBeats(t, 1, 0, 2, 1)
	checkTrucoBeats(t, 1, -1, 1, 1)
	checkTrucoBeats(t, 1, -1, -1, 0)
	checkTrucoBeats(t, 1, -1, 0, 1)

	// -1 -1 x => 0
	checkTrucoBeats(t, -1, -1, 2, 0)
	checkTrucoBeats(t, -1, 0, 2, 0)
	checkTrucoBeats(t, -1, 1, -1, 0)
	checkTrucoBeats(t, -1, 1, 1, 1)
	checkTrucoBeats(t, -1, 1, 0, 0)

	// 0 1 x => 1
	checkTrucoBeats(t, 0, 1, 2, 1)
	checkTrucoBeats(t, 0, -1, 2, 0)
	checkTrucoBeats(t, 0, 0, 1, 1)
	checkTrucoBeats(t, 0, 0, -1, 0)
	checkTrucoBeats(t, 0, 0, 0, 0)
}

func checkTrucoBeats(t *testing.T, s0, s1, s2, expected int) {
	// Helper to construct hands that produce the desired score array
	// Truco values logic:
	// mHand[i] > oHand[i] -> 1
	// mHand[i] < oHand[i] -> -1
	// mHand[i] == oHand[i] -> 0

	oHand := make(Hand, 3)
	mHand := make(Hand, 3)

	setupRound(s0, 0, &mHand, &oHand)
	setupRound(s1, 1, &mHand, &oHand)
	setupRound(s2, 2, &mHand, &oHand)

	result := TrucoBeats(mHand, oHand, NO_CARD)
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

func TestTrucoBeatsUY(t *testing.T) {
	muestra := Card{1, 'b'} // None of {1e, 4c, 1c} are piezas with 1b

	// Standard cases (no piezas)
	checkTrucoBeatsUY(t, muestra, 1, 1, 2, 1)
	checkTrucoBeatsUY(t, muestra, 1, 0, 2, 1)
	checkTrucoBeatsUY(t, muestra, 1, -1, 1, 1)
	checkTrucoBeatsUY(t, muestra, 1, -1, -1, 0)
	checkTrucoBeatsUY(t, muestra, 1, -1, 0, 1)

	checkTrucoBeatsUY(t, muestra, -1, -1, 2, 0)
	checkTrucoBeatsUY(t, muestra, -1, 0, 2, 0)
	checkTrucoBeatsUY(t, muestra, -1, 1, -1, 0)
	checkTrucoBeatsUY(t, muestra, -1, 1, 1, 1)
	checkTrucoBeatsUY(t, muestra, -1, 1, 0, 0)

	checkTrucoBeatsUY(t, muestra, 0, 1, 2, 1)
	checkTrucoBeatsUY(t, muestra, 0, -1, 2, 0)
	checkTrucoBeatsUY(t, muestra, 0, 0, 1, 1)
	checkTrucoBeatsUY(t, muestra, 0, 0, -1, 0)
	checkTrucoBeatsUY(t, muestra, 0, 0, 0, 0)

	// Pieza cases
	// muestra 1e -> 2e is pieza(19), 1e is 14
	muestraPiece := Card{1, 'e'}
	p2 := Card{2, 'e'} // 19
	c1 := Card{1, 'e'} // 14
	c4 := Card{4, 'c'} // 1

	// mHand wins with piece
	if res := TrucoBeats(Hand{p2, c4, c4}, Hand{c1, c4, c4}, muestraPiece); res != 1 {
		t.Errorf("Piece 2e(19) should beat 1e(14), got %d", res)
	}

	// oHand wins with piece
	if res := TrucoBeats(Hand{c1, c4, c4}, Hand{p2, c4, c4}, muestraPiece); res != 0 {
		t.Errorf("1e(14) should lose to 2e(19), got %d", res)
	}
}

func checkTrucoBeatsUY(t *testing.T, m Card, s0, s1, s2, expected int) {
	oHand := make(Hand, 3)
	mHand := make(Hand, 3)

	setupRound(s0, 0, &mHand, &oHand)
	setupRound(s1, 1, &mHand, &oHand)
	setupRound(s2, 2, &mHand, &oHand)

	result := TrucoBeats(mHand, oHand, m)
	if result != expected {
		t.Errorf("UY: For scores [%d, %d, %d], expected %d, got %d", s0, s1, s2, expected, result)
	}
}
