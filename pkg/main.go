package main

import (
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

	// hand := truco.Hand{{N: 7, S: 'e'}, {N: 12, S: 'c'}, {N: 2, S: 'o'}}
	// oHand := []truco.Card{{N: 3, S: 'o'}}

	// fmt.Println("soy mano, antes de jugar")
	// hand.TrucoStrengthStats([]truco.Card{}, []truco.Card{}, 255, true, true).PPrint()
	// // fmt.Println()
	// // fmt.Println("soy pie, antes de jugar")
	// // hand.TrucoStrengthStats([]truco.Card{}, []truco.Card{}, false).PPrint()

	// fmt.Println()
	// fmt.Println("el otro jugador jugo 3o")
	// hand.TrucoStrengthStats(oHand, []truco.Card{}, 255, true, true).PPrint()
	// // fmt.Println()
	// // fmt.Println("soy pie")
	// // hand.TrucoStrengthStats(oHand, []truco.Card{}, false).PPrint()

	// // TODO test HasAllInPlace and HasAll
	// // fmt.Println((truco.Hand{{11, 'o'}, {10, 'c'}, {12, 'o'}}).HasAllInPlace(truco.Hand{{12, 'o'}, {10, 'c'}}))

	// // fmt.Println(hand.TrucoStrength())

	hand := truco.Hand{{N: 6, S: 'e'}, {N: 6, S: 'c'}, {N: 6, S: 'o'}}
	for {
		hand.TrucoStrengthUY()
	}
	// hand.TrucoStrength()
}
