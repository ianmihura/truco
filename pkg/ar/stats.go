package ar

import (
	"encoding/csv"
	"fmt"
	"os"
	"runtime"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"truco/pkg/math"
)

type PairStat struct {
	Pair         string  `json:"pair"`
	IsEnvido     bool    `json:"is_envido"`
	TrucoMax     float64 `json:"truco_max"`
	TrucoMin     float64 `json:"truco_min"`
	TrucoMean    float64 `json:"truco_mean"`
	TrucoMedian  float64 `json:"truco_median"`
	EnvidoMax    int     `json:"envido_max"`
	EnvidoMin    int     `json:"envido_min"`
	EnvidoMean   float64 `json:"envido_mean"`
	EnvidoMedian float64 `json:"envido_median"`
	CombinedMean float64 `json:"combined_mean"`
	Count        int     `json:"count"`
}

// Creates a csv file that lists all possible hands with:
//   - truco strength (beatness against other hands)
//   - envido
//   - combined strength: (envido/33 + (strength+36)/72) / 2
//
// Hands are sorted by strength.
// Cards in each hand are sorted by truco strength.
func CreateHandStatsCSV(outputPath string) error {
	handsIter := math.Combinations(ALL_CARDS, 3)

	type job struct {
		hand Hand
	}
	type result struct {
		handStr  string
		strength float32
		envido   uint8
		combined float32
	}

	numWorkers := runtime.NumCPU()
	jobs := make(chan job, 100)
	resultsChan := make(chan result, 100)
	var wg sync.WaitGroup

	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				h := j.hand
				// Cards in each hand are sorted by truco strength.
				slices.SortFunc(h, SortForTruco)

				var handStr strings.Builder
				for i, c := range h {
					handStr.WriteString(c.ToString())
					if i < len(h)-1 {
						handStr.WriteString(" ")
					}
				}

				strength := (h.TrucoStrength() + 36.0) / 72
				envido := h.Envido()
				combined := (float32(envido)/MAX_ENVIDO + strength) / 2.0

				resultsChan <- result{
					handStr:  handStr.String(),
					strength: strength,
					envido:   envido,
					combined: combined,
				}
			}
		}()
	}

	// Feed jobs
	go func() {
		for h := range handsIter {
			// math.Combinations yields a slice that we should clone if we were to reuse it,
			// but its implementation already clones it.
			jobs <- job{hand: Hand(h)}
		}
		close(jobs)
	}()

	// Close results channel when workers are done
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	var allResults []result
	for r := range resultsChan {
		allResults = append(allResults, r)
	}

	// Hands are sorted by strength (descending).
	slices.SortFunc(allResults, func(a, b result) int {
		if b.strength > a.strength {
			return 1
		} else if b.strength < a.strength {
			return -1
		}
		return 0
	})

	// Save to CSV
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	// Header
	writer.Write([]string{"hand", "strength", "envido", "combined"})

	for _, r := range allResults {
		writer.Write([]string{
			r.handStr,
			fmt.Sprintf("%.6f", r.strength),
			fmt.Sprintf("%d", r.envido),
			fmt.Sprintf("%.6f", r.combined),
		})
	}

	return nil
}

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

// Reads hand strengths (output of CreateHandStatsCSV) and computes stats per pair,
// including max, min, mean, and median envido scores.
//
// This will get executed very often: every time there's a state change in the frontend.
// The frontend uses this output for stats and color coding the matrix.
func CreatePairStatsCSV(inputPath, outputPath string) error {
	records, err := getCSVReader(inputPath)
	if err != nil {
		return err
	}

	type statsKey struct {
		pair     string
		isEnvido bool
	}

	type pairData struct {
		scores  []float64
		envidos []int
	}
	statsMap := make(map[statsKey]*pairData)

	// Create stats
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
		isEnvido := envido >= 20

		// Key: first two cards (matches existing PairStrengths logic)
		r0 := NewCard(cards[0]).ToRank()
		r1 := NewCard(cards[1]).ToRank()
		pairKey := fmt.Sprintf("%s %s", r0, r1)

		key := statsKey{pair: pairKey, isEnvido: isEnvido}

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
		"is_envido",
		"truco_max",
		"truco_min",
		"truco_mean",
		"truco_median",
		"envido_max",
		"envido_min",
		"envido_mean",
		"envido_median",
		"combined_mean",
		"count",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Sort keys for stable output
	type sortedKey struct {
		pair     string
		isEnvido bool
	}
	sortedKeys := make([]sortedKey, 0, len(statsMap))
	for k := range statsMap {
		sortedKeys = append(sortedKeys, sortedKey{k.pair, k.isEnvido})
	}
	sort.Slice(sortedKeys, func(i, j int) bool {
		if sortedKeys[i].pair != sortedKeys[j].pair {
			return sortedKeys[i].pair < sortedKeys[j].pair
		}
		return sortedKeys[i].isEnvido && !sortedKeys[j].isEnvido // envido (true) first
	})

	for _, k := range sortedKeys {
		mapKey := statsKey{k.pair, k.isEnvido}
		d := statsMap[mapKey]
		if len(d.scores) == 0 {
			continue
		}

		// Truco stats
		sort.Float64s(d.scores)
		count := len(d.scores)
		minT := d.scores[0]
		maxT := d.scores[count-1]
		var sumS float64
		for _, s := range d.scores {
			sumS += s
		}
		meanT := sumS / float64(count)

		var medianT float64
		if count%2 == 1 {
			medianT = float64(d.scores[count/2])
		} else {
			medianT = float64(d.scores[count/2-1]+d.scores[count/2]) / 2.0
		}

		// Envido stats
		sort.Ints(d.envidos)
		minE := d.envidos[0]
		maxE := d.envidos[count-1]
		var sumE int
		for _, e := range d.envidos {
			sumE += e
		}
		meanE := float64(sumE) / float64(count)

		var medianE float64
		if count%2 == 1 {
			medianE = float64(d.envidos[count/2])
		} else {
			medianE = float64(d.envidos[count/2-1]+d.envidos[count/2]) / 2.0
		}

		meanC := (meanT + meanE/MAX_ENVIDO) / 2

		writer.Write([]string{
			k.pair,
			fmt.Sprintf("%v", k.isEnvido),
			fmt.Sprintf("%.6f", maxT),
			fmt.Sprintf("%.6f", minT),
			fmt.Sprintf("%.6f", meanT),
			fmt.Sprintf("%.6f", medianT),
			fmt.Sprintf("%d", maxE),
			fmt.Sprintf("%d", minE),
			fmt.Sprintf("%.2f", meanE),
			fmt.Sprintf("%.2f", medianE),
			fmt.Sprintf("%.6f", meanC),
			fmt.Sprintf("%d", count),
		})
	}

	return nil
}

func LoadPairStats(csvPath string) (map[string]PairStat, error) {
	records, err := getCSVReader(csvPath)
	if err != nil {
		return nil, err
	}

	stats := make(map[string]PairStat)
	header := records[0]

	for i := 1; i < len(records); i++ {
		row := records[i]
		if len(row) < len(header) {
			continue
		}

		data := PairStat{}
		for idx, col := range header {
			val := row[idx]
			switch col {
			case "pair":
				data.Pair = val
			case "is_envido":
				data.IsEnvido = val == "true"
			case "truco_max":
				data.TrucoMax, _ = strconv.ParseFloat(val, 64)
			case "truco_min":
				data.TrucoMin, _ = strconv.ParseFloat(val, 64)
			case "truco_mean":
				data.TrucoMean, _ = strconv.ParseFloat(val, 64)
			case "truco_median":
				data.TrucoMedian, _ = strconv.ParseFloat(val, 64)
			case "envido_max":
				v, _ := strconv.Atoi(val)
				data.EnvidoMax = v
			case "envido_min":
				v, _ := strconv.Atoi(val)
				data.EnvidoMin = v
			case "envido_mean":
				data.EnvidoMean, _ = strconv.ParseFloat(val, 64)
			case "envido_median":
				data.EnvidoMedian, _ = strconv.ParseFloat(val, 64)
			case "combined_mean":
				data.CombinedMean, _ = strconv.ParseFloat(val, 64)
			case "count":
				v, _ := strconv.Atoi(val)
				data.Count = v
			}
		}

		key := fmt.Sprintf("%s %v", data.Pair, data.IsEnvido)
		stats[key] = data
	}

	return stats, nil
}
