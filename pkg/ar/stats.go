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

// Creates a csv file that lists all possible hands with:
//   - their truco strength (beatness against other hands)
//   - envido
//   - combined strength: (envido/33 + (strength+36)/72) / 2
//
// Hands are sorted by strength.
// Cards in each hand are also sorted by truco strength.
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
				combined := (float32(envido)/33.0 + strength) / 2.0

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

// CreatePairStatsCSV reads hand strengths and computes detailed stats per pair,
// including max, min, mean, and median envido scores.
func CreatePairStatsCSV(inputPath, outputPath string) error {
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
		"envido_min_w_e",
		"envido_mean",
		"envido_mean_w_e",
		"envido_median",
		"envido_median_w_e",
		"count",
		"count_w_e",
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
		count := len(d.scores)
		if count%2 == 1 {
			medianT = float64(d.scores[count/2])
		} else {
			medianT = float64(d.scores[count/2-1]+d.scores[count/2]) / 2.0
		}

		// Envido stats
		sort.Ints(d.envidos)
		minE := d.envidos[0]
		maxE := d.envidos[len(d.envidos)-1]
		minWE := maxE // will be lower than maxE but higher than minE
		var sumE, sumWE, countWE int
		for _, e := range d.envidos {
			sumE += e
			if e >= 20 {
				sumWE += e
				countWE++
				if e < minWE {
					minWE = e
				}
			}
		}
		meanE := float64(sumE) / float64(count)
		meanWE := float64(sumWE) / float64(countWE)
		var medianE float64
		if count%2 == 1 {
			medianE = float64(d.envidos[count/2])
		} else {
			medianE = float64(d.envidos[count/2-1]+d.envidos[count/2]) / 2.0
		}

		var medianWE float64
		if countWE > 0 {
			start := count - countWE
			if countWE%2 == 1 {
				medianWE = float64(d.envidos[start+countWE/2])
			} else {
				medianWE = float64(d.envidos[start+countWE/2-1]+d.envidos[start+countWE/2]) / 2.0
			}
		}

		// TODO correct difference between envido and no envido pairs
		// probably better to do it row-wise rather than column-wise

		writer.Write([]string{
			k,
			fmt.Sprintf("%.6f", maxT),
			fmt.Sprintf("%.6f", minT),
			fmt.Sprintf("%.6f", meanT),
			fmt.Sprintf("%.6f", medianT),
			fmt.Sprintf("%d", maxE),
			fmt.Sprintf("%d", minE),
			fmt.Sprintf("%d", minWE),
			fmt.Sprintf("%.2f", meanE),
			fmt.Sprintf("%.2f", meanWE),
			fmt.Sprintf("%.1f", medianE),
			fmt.Sprintf("%.1f", medianWE),
			fmt.Sprintf("%d", count),
			fmt.Sprintf("%d", countWE),
		})
	}

	return nil
}
