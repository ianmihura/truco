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

	// fmt.Println('a', 'b', 'c', 'd', 'e')

	// fmt.Println(fsm.ValidAction("asdf"))

	hand := truco.Hand{{7, 'e'}, {2, 'c'}, {4, 'c'}}
	hand.TrucoStrengthStats([]truco.Card{{12, 'o'}}, []truco.Card{}).PPrint()
}
