package math

import (
	"iter"
	"slices"
)

// Combinations returns an iterator that yields r-length combinations of elements in the items slice.
// This replicates the behavior of python's itertools.combinations.
//
// Elements are treated as unique based on their position, not on their value.
// The combinations are emitted in lexicographic sorting order.
//
// Order does not matter in combinations, meaning (a, b) is the same as (b, a) and only one is produced.
func Combinations[T any](items []T, r int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		n := len(items)
		if r > n {
			return
		}

		indices := make([]int, r)
		for i := range indices {
			indices[i] = i
		}

		// Yield first combination
		result := make([]T, r)
		for i := range r {
			result[i] = items[indices[i]]
		}
		if !yield(slices.Clone(result)) {
			return
		}

		for {
			i := r - 1
			// Find the rightmost index i such that indices[i] != i + n - r
			for ; i >= 0; i-- {
				if indices[i] != i+n-r {
					break
				}
			}
			// If no such index exists, we are done
			if i < 0 {
				return
			}

			indices[i]++
			for j := i + 1; j < r; j++ {
				indices[j] = indices[j-1] + 1
			}

			for k := range r {
				result[k] = items[indices[k]]
			}
			if !yield(slices.Clone(result)) {
				return
			}
		}
	}
}

// Permutations returns an iterator that yields r-length permutations of elements in the items slice.
// This replicates the behavior of python's itertools.permutations.
//
// Elements are treated as unique based on their position, not on their value.
// The permutations are emitted in lexicographic sorting order (based on input order).
//
// Order matters in permutations, meaning (a, b) is different from (b, a).
func Permutations[T any](items []T, r int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		n := len(items)
		if r > n {
			return
		}

		indices := make([]int, n)
		for i := range indices {
			indices[i] = i
		}

		cycles := make([]int, r)
		for i := range cycles {
			cycles[i] = n - i
		}

		result := make([]T, r)
		for i := range r {
			result[i] = items[indices[i]]
		}
		if !yield(slices.Clone(result)) {
			return
		}

		for {
			i := r - 1
			for ; i >= 0; i-- {
				cycles[i]--
				if cycles[i] == 0 {
					tmp := indices[i]
					copy(indices[i:], indices[i+1:])
					indices[n-1] = tmp
					cycles[i] = n - i
				} else {
					j := cycles[i]
					indices[i], indices[n-j] = indices[n-j], indices[i]

					for k := range r {
						result[k] = items[indices[k]]
					}
					if !yield(slices.Clone(result)) {
						return
					}
					break
				}
			}
			if i < 0 {
				return
			}
		}
	}
}

// Hany factorial that multiplies x*...*x-l.
// Same as doing x!/l!. Note that Fact(x, 1) == x!
// Handy for a faster Pick function
func Fact(x, l int) int {
	if x <= l {
		return l
	} else {
		return x * Fact(x-1, l)
	}
}

// From n pick k
func Pick(n, k int) float32 {
	return float32(Fact(n, n-k+1)) / float32(Fact(k, 1))
}

// Type T=numeric is limited, consider adding more types as needed
type numeric interface {
	int | uint8 | float32 | float64
}

func Sum[T numeric](ar []T) T {
	var sum T
	for _, s := range ar {
		sum += s
	}
	return sum
}

func Mean[T numeric](ar []T) T {
	return Sum(ar) / T(len(ar))
}

func Median[T numeric](ar []T) T {
	count := len(ar)
	if count%2 == 1 {
		return T(ar[count/2])
	} else {
		return T(ar[count/2-1]+ar[count/2]) / 2.0
	}
}
