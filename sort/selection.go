package sort

// Selection sort uses ~N²/2 compares and N exchanges to sort any array of length N.
// So it is not suitable for sorting large arrays.
// This sorting algorithm is not stable.
// It uses Θ(1) extra memory (not including the input array).
// The complexity is O(N*N) where N = len(x).
func Selection[T any](x []T, less func(a, b T) bool) {
	n := len(x)
	for i := 0; i < n; i++ {
		min := i
		for j := i + 1; j < n; j++ {
			if less(x[j], x[min]) {
				min = j
			}
		}
		x[i], x[min] = x[min], x[i]
	}
}
