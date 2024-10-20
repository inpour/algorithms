package sort

// Merge sorts an array using a top-down, recursive version of mergesort.
// This implementation takes Θ(N*log(N)) time to sort any array of length N (assuming comparisons
// take constant time). It makes between ~N*log(N)/2 and ~N*log(N) compares.
// This sorting algorithm is stable.
// It uses Θ(N) extra memory (not including the input array).
// The complexity is O(N*log(N)) where N = len(x).
func Merge[T any](x []T, less func(a, b T) bool) {
	aux := make([]T, len(x))
	sortMerge(x, aux, 0, len(x)-1, less)
}

func sortMerge[T any](x []T, aux []T, lo, hi int, less func(a, b T) bool) {
	if hi <= lo {
		return
	}
	mid := lo + (hi-lo)/2
	sortMerge(x, aux, lo, mid, less)
	sortMerge(x, aux, mid+1, hi, less)
	merge(x, aux, lo, mid, hi, less)
}

func merge[T any](x []T, aux []T, lo, mid, hi int, less func(a, b T) bool) {
	// copy x to aux
	for k := lo; k <= hi; k++ {
		aux[k] = x[k]
	}

	// merge back to x
	i := lo
	j := mid + 1
	for k := lo; k <= hi; k++ {
		if i > mid {
			x[k] = aux[j]
			j++
		} else if j > hi {
			x[k] = aux[i]
			i++
		} else if less(aux[j], aux[i]) {
			x[k] = aux[j]
			j++
		} else {
			x[k] = aux[i]
			i++
		}
	}
}
