package sort

// Heap takes ~Nlog(N) time to sort any array of length N (assuming comparisons take constant time).
// It makes at most 2N*log(N) compares.
// This sorting algorithm is not stable.
// It uses ~1 extra memory (not including the input array).
// The complexity is O(N*log(N)) where N = len(x).
func Heap[T any](x []T, less func(a, b T) bool) {
	n := len(x)

	// heapify
	for k := n / 2; k >= 1; k-- {
		sink(x, k, n, less)
	}

	// sort-down
	k := n
	for k > 1 {
		x[0], x[k-1] = x[k-1], x[0]
		k--
		sink[T](x, 1, k, less)
	}
}

func sink[T any](x []T, k, n int, less func(a, b T) bool) {
	for 2*k <= n {
		j := 2 * k
		if j < n && less(x[j-1], x[j]) {
			j++
		}
		if !less(x[k-1], x[j-1]) {
			break
		}
		x[k-1], x[j-1] = x[j-1], x[k-1]
		k = j
	}
}
