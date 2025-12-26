package main

import (
	"truco/pkg/fsm"
)

func main() {
	// if err := ar.CreateHandStatsCSV("web/static/hand_stats.csv"); err != nil {
	// 	fmt.Println("Error generating hand strength CSV:", err)
	// } else {
	// 	fmt.Println("Successfully generated hand strength CSV in web/static/hand_stats.csv")
	// }

	// if err := ar.CreatePairStatsCSV("web/static/hand_stats.csv", "web/static/pair_stats.csv"); err != nil {
	// 	fmt.Println("Error generating pair stats:", err)
	// } else {
	// 	fmt.Println("Successfully generated pair stats with envido info in web/static/pair_stats.csv")
	// }

	fsm.NewMatch()
}
