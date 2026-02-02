package truco

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
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

type StatsKey struct {
	pair     string
	isEnvido bool
}

type PairData struct {
	scores  []float64
	envidos []int
}

type FilterHands struct {
	KCards  []Card
	MCards  []Card
	KEnvido []uint8 // TODO maybe unnecesary
	MEnvido uint8
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
				combined := (float32(envido)/MAX_ENVIDO_AR + strength) / 2.0

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

	// Read and discard the header row
	if _, err := reader.Read(); err != nil {
		if err == io.EOF {
			log.Fatalf("Empty file or file with only a header")
		} else {
			log.Fatalf("Error reading header: %v", err)
		}
	}

	r, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func filterRecords(records *[][]string, filter FilterHands) *[][]string {
	if len(filter.KCards) == 0 && len(filter.MCards) == 0 && filter.MEnvido == 255 {
		return records
	}

	// Create a map of known cards for efficient lookup
	kCardsSet := make(map[string]struct{}, len(filter.KCards))
	for _, c := range filter.KCards {
		kCardsSet[c.ToString()] = struct{}{}
	}

	fRecords := make([][]string, 0, len(*records))
	for _, row := range *records {
		if len(row) == 0 {
			continue
		}

		handStr := row[0]
		cards := strings.Split(handStr, " ")

		// Known cards filter: hand should NOT contain any kCard
		hasKCard := false
		for _, cStr := range cards {
			if _, found := kCardsSet[cStr]; found {
				hasKCard = true
				break
			}
		}
		if hasKCard {
			continue
		}

		// My cards filter: hand MUST contain all mCards
		hasAllMCards := true
		for _, mCard := range filter.MCards {
			mCardStr := mCard.ToString()
			found := slices.Contains(cards, mCardStr)
			if !found {
				hasAllMCards = false
				break
			}
		}
		if !hasAllMCards {
			continue
		}

		if filter.MEnvido < 100 {
			// MEnvido declared exactly
			if NewHand(handStr).Envido() != filter.MEnvido {
				continue
			}
		} else if filter.MEnvido != 255 {
			// MEnvido declared at a range (eg. 127: '27 son buenas')
			// This means my hand is worse than or equal to 27.
			if NewHand(handStr).Envido() > filter.MEnvido-100 {
				continue
			}
		}

		fRecords = append(fRecords, row)
	}

	return &fRecords
}

// Reads hand strengths (output of CreateHandStatsCSV) and computes stats per pair,
// including max, min, mean, and median scores.
//
// Input params:
//   - withEnvido: if true, discriminates between hands with and without envido (PK = pair + is_envido)
//   - filter: filter out impossible hands
//
// This is executed on every state change in the tracker to provide real-time hand strength feedback.
func ComputePairStats(withEnvido bool, filter FilterHands) (map[string]PairStat, error) {
	records, err := getCSVReader("web/static/hand_stats.csv")
	if err != nil {
		return nil, err
	}
	records = *filterRecords(&records, filter)

	statsMapInternal := make(map[StatsKey]*PairData)

	// Ingest records into internal map
	for _, row := range records {
		handStr := row[0]
		scoreStr := row[1]
		envidoStr := row[2]

		score, err := strconv.ParseFloat(scoreStr, 64)
		if err != nil {
			continue
		}

		envido, err := strconv.ParseInt(envidoStr, 10, 8)
		if err != nil {
			continue
		}
		isEnvido := envido >= 20

		cards := strings.Split(handStr, " ")
		if len(cards) < 3 {
			continue
		}

		// Key: first two cards (assumes cards are already sorted for truco)
		r0 := NewCard(cards[0]).ToRank()
		r1 := NewCard(cards[1]).ToRank()
		pairKey := fmt.Sprintf("%s %s", r0, r1)

		// Full hand for envido
		var key StatsKey
		if withEnvido {
			key = StatsKey{pair: pairKey, isEnvido: isEnvido}
		} else {
			key = StatsKey{pair: pairKey, isEnvido: false}
		}

		if statsMapInternal[key] == nil {
			statsMapInternal[key] = &PairData{}
		}
		statsMapInternal[key].scores = append(statsMapInternal[key].scores, score)
		statsMapInternal[key].envidos = append(statsMapInternal[key].envidos, int(envido))
	}

	statsResult := make(map[string]PairStat)

	// Compute metrics for each pair
	for key, d := range statsMapInternal {
		if len(d.scores) == 0 {
			continue
		}

		// Truco stats
		sort.Float64s(d.scores)
		count := len(d.scores)
		minT := d.scores[0]
		maxT := d.scores[count-1]
		meanT := math.Mean(d.scores)
		medianT := math.Median(d.scores)

		// Envido stats
		sort.Ints(d.envidos)
		minE := d.envidos[0]
		maxE := d.envidos[count-1]
		sumE := math.Sum(d.envidos)
		meanE := float64(sumE) / float64(count)
		medianE := math.Median(d.envidos)

		meanC := (meanT + meanE/MAX_ENVIDO_AR) / 2

		stat := PairStat{
			Pair:         key.pair,
			IsEnvido:     key.isEnvido,
			TrucoMax:     maxT,
			TrucoMin:     minT,
			TrucoMean:    meanT,
			TrucoMedian:  medianT,
			EnvidoMax:    maxE,
			EnvidoMin:    minE,
			EnvidoMean:   meanE,
			EnvidoMedian: float64(medianE),
			CombinedMean: meanC,
			Count:        count,
		}

		// Return key matches frontend expectations: "rank1 rank2 bool"
		mapKey := fmt.Sprintf("%s %v", stat.Pair, stat.IsEnvido)
		statsResult[mapKey] = stat
	}

	return statsResult, nil
}
