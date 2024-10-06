package sort

// Insertion uses ~N²/4 compares and ~N²/4 exchanges to sort a randomly ordered array
// of length N with distinct keys, on the average. The worst case is ~N²/2 compares and
// ~N²/2 exchanges and the best cas is N-1 compares and 0 exchanges.
// So, it is not suitable for sorting large arbitrary arrays. More precisely, the number
// of exchanges is exactly equal to the number of inversions. For example, it sorts a
// partially-sorted array in linear time.
// Insertion sort is an efficient method for partially-sorted arrays. Indeed, when the
// number of inversions is low, insertion sort is likely to be faster than any sorting
// algorithm.
// This sorting algorithm is stable.
// It uses Θ(1) extra memory (not including the input array).
// The complexity is O(N*N) where N = len(x).
func Insertion[T any](x []T, less func(a, b T) bool) {
	n := len(x)
	for i := 1; i < n; i++ {
		for j := i; j > 0 && less(x[j], x[j-1]); j-- {
			x[j], x[j-1] = x[j-1], x[j]
		}
	}
}
