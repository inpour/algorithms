package sort

// Shell sorts an array with Knuth's increment sequence (https://oeis.org/A003462).
// In the worst case, this implementation makes Θ(N^(3/2)) compares and exchanges to sort
// an array of length N.
// This sorting algorithm is not stable.
// It uses Θ(1) extra memory (not including the input array).
// The complexity is not precisely characterized! Maybe N*log(N) or N^(6/5), where N = len(x).
func Shell[T any](x []T, less func(a, b T) bool) {
	n := len(x)
	// 3x+1 increment sequence:  1, 4, 13, 40, 121, 364, 1093, ...
	h := 1
	for h < n/3 {
		h = 3*h + 1
	}
	for h >= 1 {
		// h-sort the array
		for i := h; i < n; i++ {
			for j := i; j >= h && less(x[j], x[j-h]); j -= h {
				x[j], x[j-h] = x[j-h], x[j]
			}
		}
		h /= 3
	}
}
