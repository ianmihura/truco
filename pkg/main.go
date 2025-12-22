package main

import (
	"fmt"
	"truco/pkg/ar"
	"truco/pkg/math"
)

func main() {
	// h := Hand{Card{3, 'e'}, Card{4, 'e'}, Card{3, 'b'}}
	// h.Println()

	// r := CardRange(127, []Card{{10, 'c'}, {4, 'e'}}, []Card{{10, 'o'}, {5, 'e'}})
	// for i := range r {
	// 	r[i].Println()
	// }

	// h := []Card{{10, 'c'}, {4, 'e'}, {1, 'e'}, {7, 'o'}, {1, 'b'}, {2, 'o'}}
	// slices.SortFunc(h, sortForTruco)
	// fmt.Println(h)

	// var mHand, oHand Hand
	// mHand = Hand{{2, 'b'}, {3, 'o'}, {4, 'b'}}
	// oHand = Hand{{7, 'b'}, {7, 'o'}, {4, 'o'}}
	// fmt.Println(mHand.TrucoBeatsAll(oHand))

	// mHand = Hand{{1, 'e'}, {6, 'o'}, {7, 'e'}}
	// mHand = Hand{{1, 'e'}, {1, 'b'}, {7, 'e'}}
	// mHand = Hand{{6, 'e'}, {7, 'b'}, {6, 'o'}}
	// fmt.Println(mHand.TrucoStrength())

	hands := math.Combinations(ar.ALL_CARDS, 3)
	var mHand ar.Hand
	for h := range hands {
		mHand = ar.Hand(h)
		mHand.Print()
		fmt.Print(" scores: ", mHand.TrucoStrength(), "\n")
	}
}
