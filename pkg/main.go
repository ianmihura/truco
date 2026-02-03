package main

import (
	"fmt"
	"truco/pkg/truco"
)

func main() {
	// if err := truco.CreateHandStatsCSV("web/static/hand_stats.csv"); err != nil {
	// 	fmt.Println("Error generating hand strength CSV:", err)
	// } else {
	// 	fmt.Println("Successfully generated hand strength CSV in web/static/hand_stats.csv")
	// }

	// if err := truco.CreatePairStatsCSV("web/static/hand_stats.csv", "web/static/pair_stats.csv"); err != nil {
	// 	fmt.Println("Error generating pair stats:", err)
	// } else {
	// 	fmt.Println("Successfully generated pair stats with envido info in web/static/pair_stats.csv")
	// }

	// fsm.NewMatch()

	// TODO see what happens if you play optimally (as script suggests)

	hand := truco.Hand{{7, 'e'}, {12, 'c'}, {4, 'o'}}
	oHand := []truco.Card{{12, 'o'}}

	fmt.Println("soy mano, antes de jugar")
	hand.TrucoStrengthStats([]truco.Card{}, []truco.Card{}, true).PPrint()
	fmt.Println()
	fmt.Println("soy pie, antes de jugar")
	hand.TrucoStrengthStats([]truco.Card{}, []truco.Card{}, false).PPrint()

	fmt.Println()
	fmt.Println("soy mano")
	hand.TrucoStrengthStats(oHand, []truco.Card{}, true).PPrint()
	fmt.Println()
	fmt.Println("soy pie")
	hand.TrucoStrengthStats(oHand, []truco.Card{}, false).PPrint()

	// TODO test HasAllInPlace and HasAll
	// fmt.Println((truco.Hand{{11, 'o'}, {10, 'c'}, {12, 'o'}}).HasAllInPlace(truco.Hand{{12, 'o'}, {10, 'c'}}))

	// fmt.Println(hand.TrucoStrength())
}
