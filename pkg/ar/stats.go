package ar

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func getCSVReader(csvPath string) ([][]string, error) {
	f, err := os.Open(csvPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)

	r, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func EnrichWithEnvido(csvPath, outputPath string) error {
	// Read input CSV
	f, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Create output CSV
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	// Process records
	for i, row := range records {
		// Handle header row
		if i == 0 && (strings.HasPrefix(row[0], "hand") || strings.Contains(row[0], "score")) {
			// Add "envido" column to header
			newHeader := append(row, "envido")
			if err := writer.Write(newHeader); err != nil {
				return err
			}
			continue
		}

		// Parse hand and compute envido
		if len(row) < 1 {
			continue
		}

		cards := strings.Split(row[0], " ")
		if len(cards) < 3 {
			continue
		}

		hand := Hand{
			NewCard(cards[0]),
			NewCard(cards[1]),
			NewCard(cards[2]),
		}

		envidoScore := hand.Envido()

		// Append envido score to the row
		newRow := append(row, fmt.Sprintf("%d", envidoScore))
		if err := writer.Write(newRow); err != nil {
			return err
		}
	}

	return nil
}

// PairStatsToCSV reads hand strengths and computes detailed stats per pair,
// including max, min, mean, and median envido scores.
func PairStatsToCSV(inputPath, outputPath string) error {
	records, err := getCSVReader(inputPath)
	if err != nil {
		return err
	}

	type pairData struct {
		scores  []float64
		envidos []int
	}
	statsMap := make(map[string]*pairData)

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
		if len(cards) < 3 {
			continue
		}

		// Full hand for envido
		h := Hand{NewCard(cards[0]), NewCard(cards[1]), NewCard(cards[2])}
		envido := int(h.Envido())

		// Key: first two cards (matches existing PairStrengths logic)
		r0 := NewCard(cards[0]).ToRank()
		r1 := NewCard(cards[1]).ToRank()
		key := fmt.Sprintf("%s %s", r0, r1)

		if statsMap[key] == nil {
			statsMap[key] = &pairData{}
		}
		statsMap[key].scores = append(statsMap[key].scores, score)
		statsMap[key].envidos = append(statsMap[key].envidos, envido)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	// Write header
	header := []string{
		"pair",
		"truco_max",
		"truco_min",
		"truco_mean",
		"truco_median",
		"envido_max",
		"envido_min",
		"envido_mean",
		"envido_median",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Sort keys for stable output
	keys := make([]string, 0, len(statsMap))
	for k := range statsMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		d := statsMap[k]
		if len(d.scores) == 0 {
			continue
		}

		// Truco stats
		sort.Float64s(d.scores)
		minT := d.scores[0]
		maxT := d.scores[len(d.scores)-1]
		var sumS float64
		for _, s := range d.scores {
			sumS += s
		}
		meanT := sumS / float64(len(d.scores))

		var medianT float64
		n := len(d.scores)
		if n%2 == 1 {
			medianT = float64(d.scores[n/2])
		} else {
			medianT = float64(d.scores[n/2-1]+d.scores[n/2]) / 2.0
		}

		// Envido stats
		sort.Ints(d.envidos)
		minE := d.envidos[0]
		maxE := d.envidos[len(d.envidos)-1]
		var sumE int
		for _, e := range d.envidos {
			sumE += e
		}
		meanE := float64(sumE) / float64(len(d.envidos))

		var medianE float64
		n = len(d.envidos)
		if n%2 == 1 {
			medianE = float64(d.envidos[n/2])
		} else {
			medianE = float64(d.envidos[n/2-1]+d.envidos[n/2]) / 2.0
		}

		writer.Write([]string{
			k,
			fmt.Sprintf("%.6f", maxT),
			fmt.Sprintf("%.6f", minT),
			fmt.Sprintf("%.6f", meanT),
			fmt.Sprintf("%.6f", medianT),
			fmt.Sprintf("%d", maxE),
			fmt.Sprintf("%d", minE),
			fmt.Sprintf("%.2f", meanE),
			fmt.Sprintf("%.1f", medianE),
		})
	}

	return nil
}
