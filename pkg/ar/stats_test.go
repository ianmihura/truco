package ar

import (
	"reflect"
	"testing"
)

func TestFilterRecords(t *testing.T) {
	records := [][]string{
		{"1e 7e 4b", "0.9", "28"},   // Hand 0: 1e, 7e, 4b. Envido: 1+7+20 = 28
		{"1e 2e 3e", "1.0", "25"},   // Hand 1: 1e, 2e, 3e. Envido: 2+3+20 = 25
		{"1b 2b 3c", "0.5", "23"},   // Hand 2: 1b, 2b, 3c. Envido: 1+2+20 = 23
		{"4c 5c 6c", "0.1", "31"},   // Hand 3: 4c, 5c, 6c. Envido: 5+6+20 = 31
		{"7b 10b 11b", "0.2", "27"}, // Hand 4: 7b, 10b, 11b. Envido: 7+0+20 = 27
	}

	tests := []struct {
		name     string
		filter   FilterHands
		expected [][]string
	}{
		{
			name: "Filter KCards - exclude 1e",
			filter: FilterHands{
				KCards:  []Card{{1, 'e'}},
				MEnvido: 255,
			},
			expected: [][]string{
				{"1b 2b 3c", "0.5", "23"},
				{"4c 5c 6c", "0.1", "31"},
				{"7b 10b 11b", "0.2", "27"},
			},
		},
		{
			name: "Filter MCards - must have 1e, exclude 4b",
			filter: FilterHands{
				KCards:  []Card{{4, 'b'}},
				MCards:  []Card{{1, 'e'}},
				MEnvido: 255,
			},
			expected: [][]string{
				{"1e 2e 3e", "1.0", "25"},
			},
		},
		{
			name: "Filter MEnvido - exactly 27",
			filter: FilterHands{
				KCards:  []Card{},
				MCards:  []Card{},
				MEnvido: 27,
			},
			expected: [][]string{
				{"7b 10b 11b", "0.2", "27"},
			},
		},
		{
			name: "Filter MCards and KCards",
			filter: FilterHands{
				KCards:  []Card{{7, 'e'}}, // Exclude hand 0
				MCards:  []Card{{1, 'e'}}, // Must have 1e (Hands 0, 1)
				MEnvido: 255,
			},
			expected: [][]string{
				{"1e 2e 3e", "1.0", "25"},
			},
		},
		{
			name:   "Filter MEnvido range (son buenas) - eg 27 son buenas (127)",
			filter: FilterHands{MEnvido: 127}, // 27 son buenas
			expected: [][]string{
				{"1e 2e 3e", "1.0", "25"},   // Envido 25 <= 27
				{"1b 2b 3c", "0.5", "23"},   // Envido 23 <= 27
				{"7b 10b 11b", "0.2", "27"}, // Envido 27 <= 27
			},
		},
		{
			name: "Exclude all",
			filter: FilterHands{
				KCards:  []Card{{4, 'b'}},
				MCards:  []Card{{1, 'e'}},
				MEnvido: 124,
			},
			expected: [][]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterRecords(&records, tt.filter)
			if !reflect.DeepEqual(*got, tt.expected) {
				t.Errorf("filterRecords() = %v, want %v", *got, tt.expected)
			}
		})
	}
}

func TestFilterRecordsMoreCases(t *testing.T) {
	records := [][]string{
		{"1e 7e 4b", "0.9", "28"},  // Hand 0: 1e, 7e, 4b. Envido: 1+7+20 = 28
		{"1e 2e 3e", "1.0", "25"},  // Hand 1: 1e, 2e, 3e. Envido: 2+3+20 = 25
		{"1b 2b 3c", "0.5", "23"},  // Hand 2: 1b, 2b, 3c. Envido: 1+2+20 = 23
		{"3c 5c 6c", "0.1", "31"},  // Hand 3: 4c, 5c, 6c. Envido: 5+6+20 = 31
		{"7b 3c 11b", "0.2", "27"}, // Hand 4: 7b, 10b, 11b. Envido: 7+0+20 = 27
		{"7e 3c 11b", "0.2", "7"},
	}

	tests := []struct {
		name     string
		filter   FilterHands
		expected [][]string
	}{
		{
			name: "Filter MCards - must have 1e, exclude 4b",
			filter: FilterHands{
				KCards:  []Card{{11, 'b'}},
				MCards:  []Card{{3, 'c'}},
				MEnvido: 130,
			},
			expected: [][]string{
				{"1b 2b 3c", "0.5", "23"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterRecords(&records, tt.filter)
			if !reflect.DeepEqual(*got, tt.expected) {
				t.Errorf("filterRecords() = %v, want %v", *got, tt.expected)
			}
		})
	}
}
