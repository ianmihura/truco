package math

import (
	"slices"
	"testing"
)

func TestCombinations(t *testing.T) {
	items := []string{"A", "B", "C"}
	r := 2
	expected := [][]string{
		{"A", "B"}, {"A", "C"},
		{"B", "C"},
	}

	var got [][]string
	for p := range Combinations(items, r) {
		got = append(got, p)
	}

	if len(got) != len(expected) {
		t.Fatalf("Expected %d combinations, got %d", len(expected), len(got))
	}

	for i, p := range got {
		if !slices.Equal(p, expected[i]) {
			t.Errorf("Index %d: expected %v, got %v", i, expected[i], p)
		}
	}
}

func TestCombinationsOrderDoesNotMatter(t *testing.T) {
	items := []int{1, 2}
	r := 2
	// For combinations, (1, 2) is present but (2, 1) is NOT (as it is considered the same).
	found12 := false
	found21 := false

	count := 0
	for p := range Combinations(items, r) {
		count++
		if p[0] == 1 && p[1] == 2 {
			found12 = true
		}
		if p[0] == 2 && p[1] == 1 {
			found21 = true
		}
	}

	if count != 1 {
		t.Errorf("Expected 1 combination for size 2 from 2 items, got %d", count)
	}
	if !found12 {
		t.Errorf("Expected (1,2) to be present")
	}
	if found21 {
		t.Errorf("Did not expect (2,1) to be present")
	}
}

func TestPermutations(t *testing.T) {
	items := []string{"A", "B", "C"}
	r := 2
	expected := [][]string{
		{"A", "B"}, {"A", "C"},
		{"B", "A"}, {"B", "C"},
		{"C", "A"}, {"C", "B"},
	}

	var got [][]string
	for p := range Permutations(items, r) {
		got = append(got, p)
	}

	if len(got) != len(expected) {
		t.Fatalf("Expected %d permutations, got %d", len(expected), len(got))
	}

	for i, p := range got {
		if !slices.Equal(p, expected[i]) {
			t.Errorf("Index %d: expected %v, got %v", i, expected[i], p)
		}
	}
}

func TestPermutationsOrderMatters(t *testing.T) {
	items := []int{1, 2}
	r := 2
	expected := [][]int{
		{1, 2},
		{2, 1},
	}

	var got [][]int
	for p := range Permutations(items, r) {
		got = append(got, p)
	}

	if len(got) != len(expected) {
		t.Fatalf("Expected %d permutations, got %d", len(expected), len(got))
	}

	for i, p := range got {
		if !slices.Equal(p, expected[i]) {
			t.Errorf("Index %d: expected %v, got %v", i, expected[i], p)
		}
	}
}

func TestPermutationsEmpty(t *testing.T) {
	items := []int{1, 2, 3}
	r := 4
	count := 0
	for range Permutations(items, r) {
		count++
	}
	if count != 0 {
		t.Errorf("Expected 0 permutations when r > len(items), got %d", count)
	}
}
