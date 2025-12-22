package ar

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Returns the stats map: Key="Rank1 Rank2" (Top 2 ranks sorted), Value=Average Strength
func PairStrengths(csvPath string) (map[string]float64, error) {
	f, err := os.Open(csvPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	// Skip header if present, or check first line. Logic assumes header "hand,truco_score"
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Key: "Rank1 Rank2", Value: Sum of scores
	sums := make(map[string]float64)
	counts := make(map[string]int)

	for i, row := range records {
		if i == 0 && (strings.HasPrefix(row[0], "hand") || strings.Contains(row[0], "score")) {
			continue // Skip Header
		}
		if len(row) < 2 {
			continue
		}

		handStr := row[0]
		scoreStr := row[1]

		score, err := strconv.ParseFloat(scoreStr, 64)
		if err != nil {
			continue
		}

		cards := strings.Split(handStr, " ")
		if len(cards) < 2 {
			continue
		}

		// key eg. "1e 1b"
		r0 := NewCard(cards[0]).ToRank()
		r1 := NewCard(cards[1]).ToRank()
		key := fmt.Sprintf("%s %s", r0, r1)

		sums[key] += score
		counts[key]++
	}

	averages := make(map[string]float64)
	for k, sum := range sums {
		averages[k] = sum / float64(counts[k])
	}
	return averages, nil
}
