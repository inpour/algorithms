package sort

import (
	"math/rand"
	"time"
)

// Quick sorts an array using quicksort. Quicksort is popular because it is not difficult to implement,
// works well for a variety of different kinds of input data, and is substantially faster than any
// other sorting method in typical applications.
// This implementation uses ~2N*ln(N) compares (and one-sixth that many exchanges) on the average to
// sort an array of length N with distinct keys. Quicksort uses ~N²/2 compares in the worst case,
// but random shuffling protects against this case.
// This sorting algorithm is not stable.
// It is in-place (uses only a small auxiliary stack), requires time proportional to N*log(N) on
// the average to sort N items.
// The complexity is O(N*log(N)) where N = len(x).
func Quick[T any](x []T, less func(a, b T) bool) {
	// Shuffle x
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(x), func(i, j int) {
		x[i], x[j] = x[j], x[i]
	})

	sortQuick(x, 0, len(x)-1, less)
}

func sortQuick[T any](x []T, lo, hi int, less func(a, b T) bool) {
	if hi <= lo {
		return
	}
	j := partition(x, lo, hi, less)
	sortQuick(x, lo, j-1, less)
	sortQuick(x, j+1, hi, less)
}

func partition[T any](x []T, lo, hi int, less func(a, b T) bool) int {
	i := lo
	j := hi + 1
	v := x[lo]
	for true {
		// find item on lo to swap
		for i++; less(x[i], v); i++ {
			if i == hi {
				break
			}
		}
		// find item on hi to swap
		for j--; less(v, x[j]); j-- {
			if j == lo {
				break // redundant since x[lo] acts as sentinel
			}
		}
		// check if pointers cross
		if i >= j {
			break
		}

		x[i], x[j] = x[j], x[i]
	}

	// put partitioning item v at x[j]
	x[lo], x[j] = x[j], x[lo]

	// now, x[lo .. j-1] <= x[j] <= x[j+1 .. hi]
	return j
}
