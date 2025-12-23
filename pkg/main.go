package main

import (
	"fmt"
	"truco/pkg/ar"
)

func main() {

	// hands := math.Combinations(ar.ALL_CARDS, 3)
	// var mHand ar.Hand
	// for h := range hands {
	// 	mHand = ar.Hand(h)
	// 	mHand.Print()
	// 	fmt.Print(" scores: ", mHand.TrucoStrength(), "\n")
	// }

	// TODO integrate stats into a concurrent pipeline to generate hand_strength.csv

	if err := ar.PairStatsToCSV("web/static/hand_strength.csv", "web/static/pair_strength.csv"); err != nil {
		fmt.Println("Error generating pair stats:", err)
	} else {
		fmt.Println("Successfully generated pair stats with envido info in web/static/pair_strength.csv")
	}
}
