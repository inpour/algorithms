package search

import (
	"iter"
)

// BinarySearchST represents an ordered symbol table of generic key-value pairs.
// It relies on the compare() function to compare two keys:
//
//	if a == b then compare(a, b) returns 0
//	if a > b then compare(a, b) returns 1
//	if a < b then compare(a, b) returns -1
//
// This implementation uses a sorted array and binary search.
type BinarySearchST[K, V any] struct {
	n       int              // number of key-value pairs
	keys    []K              // underlying array represents the keys
	vals    []V              // underlying array represents the vals
	compare func(a, b K) int // function to compare two keys
}

// NewBinarySearchST initializes an empty symbol table with the specified initial capacity.
// It gets a function as a parameter to compare two keys.
// The complexity is O(1).
func NewBinarySearchST[K, V any](capacity int, compare func(a, b K) int) *BinarySearchST[K, V] {
	return &BinarySearchST[K, V]{
		n:       0,
		keys:    make([]K, capacity),
		vals:    make([]V, capacity),
		compare: compare,
	}
}

// resize the underlying arrays
// The complexity is O(N) where N is the number of key-value pairs.
func (b *BinarySearchST[K, V]) resize(capacity int) {
	if capacity < b.n {
		return
	}
	tempKeys := make([]K, capacity)
	tempVals := make([]V, capacity)
	for i := 0; i < b.n; i++ {
		tempKeys[i] = b.keys[i]
		tempVals[i] = b.vals[i]
	}
	b.keys = tempKeys
	b.vals = tempVals
}

// Size returns the number of key-value pairs.
// The complexity is O(1).
func (b *BinarySearchST[K, V]) Size() int {
	return b.n
}

// IsEmpty returns true if this symbol table is empty.
// The complexity is O(1).
func (b *BinarySearchST[K, V]) IsEmpty() bool {
	return b.n == 0
}

// Contains returns true if this symbol table contain the given key.
// The complexity is O(log(N)) where N is the number of key-value pairs.
func (b *BinarySearchST[K, V]) Contains(key K) bool {
	_, err := b.Rank(key)
	return err == nil
}

// Get returns the value associated with the given key, ErrAbsentKey if key is absent.
// The complexity is O(log(N)) where N is the number of key-value pairs.
func (b *BinarySearchST[K, V]) Get(key K) (V, error) {
	r, err := b.Rank(key)
	if err != nil {
		var val V
		return val, err
	}
	return b.vals[r], nil
}

// Rank returns the number of keys strictly less than key, ErrAbsentKey if key is absent.
// The complexity is O(log(N)) where N is the number of key-value pairs.
func (b *BinarySearchST[K, V]) Rank(key K) (int, error) {
	lo := 0
	hi := b.n - 1
	for lo <= hi {
		mid := (lo + hi) / 2
		cmp := b.compare(key, b.keys[mid])
		if cmp < 0 {
			hi = mid - 1
		} else if cmp > 0 {
			lo = mid + 1
		} else {
			return mid, nil
		}
	}
	return lo, ErrAbsentKey
}

// Put inserts the specified key-value pair into the symbol table, overwriting the old value with the
// new value if the symbol table already contains the specified key.
// The complexity is O(N) where N is the number of key-value pairs.
func (b *BinarySearchST[K, V]) Put(key K, val V) {
	i, err := b.Rank(key)

	// key is already in table
	if err == nil {
		b.vals[i] = val
		return
	}

	// insert new key-value pair
	if b.n == cap(b.keys) {
		b.resize(2 * b.n)
	}
	for j := b.n; j > i; j-- {
		b.keys[j] = b.keys[j-1]
		b.vals[j] = b.vals[j-1]
	}
	b.keys[i] = key
	b.vals[i] = val
	b.n++
}

// Delete removes the specified key and associated value, ErrAbsentKey if key is absent.
// The complexity is O(N) where N is the number of key-value pairs.
func (b *BinarySearchST[K, V]) Delete(key K) error {
	i, err := b.Rank(key)

	// key not in table
	if err != nil {
		return ErrAbsentKey
	}

	for j := i; j < b.n-1; j++ {
		b.keys[j] = b.keys[j+1]
		b.vals[j] = b.vals[j+1]
	}
	b.n--

	// resize if 1/4 full
	if b.n > 0 && b.n == cap(b.keys)/4 {
		b.resize(2 * b.n)
	}
	return nil
}

// DelMin removes the smallest key and associated value, ErrEmptySymbolTable if the symbol table is empty.
// The complexity is O(N) where N is the number of key-value pairs.
func (b *BinarySearchST[K, V]) DelMin() error {
	key, err := b.Min()
	if err != nil {
		return err
	}
	b.Delete(key)
	return nil
}

// DelMax removes the largest key and associated value, ErrEmptySymbolTable if the symbol table is empty.
// The complexity is O(N) where N is the number of key-value pairs.
func (b *BinarySearchST[K, V]) DelMax() error {
	key, err := b.Max()
	if err != nil {
		return err
	}
	b.Delete(key)
	return nil
}

// Min returns the smallest key, ErrEmptySymbolTable if the symbol table is empty.
// The complexity is O(1).
func (b *BinarySearchST[K, V]) Min() (K, error) {
	var key K
	if b.IsEmpty() {
		return key, ErrEmptySymbolTable
	}
	return b.keys[0], nil
}

// Max returns the largest key, ErrEmptySymbolTable if the symbol table is empty.
// The complexity is O(1).
func (b *BinarySearchST[K, V]) Max() (K, error) {
	var key K
	if b.IsEmpty() {
		return key, ErrEmptySymbolTable
	}
	return b.keys[b.n-1], nil
}

// Select return the kth smallest key (key of rank k), ErrInvalidRank if rank is out of range.
// The complexity is O(1).
func (b *BinarySearchST[K, V]) Select(k int) (K, error) {
	var key K
	if k < 0 || k >= b.n {
		return key, ErrInvalidRank
	}
	return b.keys[k], nil
}

// Floor returns the largest key less than or equal to key, ErrTooSmallFloorKey if key to floor is too small.
// The complexity is O(log(N)) where N is the number of key-value pairs.
func (b *BinarySearchST[K, V]) Floor(key K) (K, error) {
	i, err := b.Rank(key)
	if err == nil {
		return b.keys[i], nil
	}
	if i == 0 {
		return key, ErrTooSmallFloorKey
	}
	return b.keys[i-1], nil
}

// Ceiling returns the smallest key greater than or equal to key, ErrTooLargeCeilingKey if key to ceiling is too large.
// The complexity is O(log(N)) where N is the number of key-value pairs.
func (b *BinarySearchST[K, V]) Ceiling(key K) (K, error) {
	i, _ := b.Rank(key)
	if i == b.n {
		return key, ErrTooLargeCeilingKey
	}
	return b.keys[i], nil
}

// RangeSize returns the number of keys in [lo:hi] range.
// The complexity is O(log(N)) where N is the number of key-value pairs.
func (b *BinarySearchST[K, V]) RangeSize(lo, hi K) int {
	if b.compare(lo, hi) > 0 {
		return 0
	}
	hiRank, err := b.Rank(hi)
	loRank, _ := b.Rank(lo)
	size := hiRank - loRank
	// symbol table contains hi key
	if err == nil {
		size++
	}
	return size
}

// Iterator returns an iterator that iterates over all key-value pairs in sorted order.
func (b *BinarySearchST[K, V]) Iterator() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for i := 0; i < b.n; i++ {
			if !yield(b.keys[i], b.vals[i]) {
				return
			}
		}
	}
}

// RangeIterator returns an iterator that iterates over key-value pairs where keys in [lo:hi] range, in sorted order.
// It takes O(log(N)) time to prepare iterator where N is the number of key-value pairs.
func (b *BinarySearchST[K, V]) RangeIterator(lo, hi K) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if b.compare(lo, hi) > 0 {
			return
		}
		hiRank, err := b.Rank(hi)
		for i, _ := b.Rank(lo); i < hiRank; i++ {
			if !yield(b.keys[i], b.vals[i]) {
				return
			}
		}
		// symbol table contains hi key
		if err == nil {
			if !yield(b.keys[hiRank], b.vals[hiRank]) {
				return
			}
		}
	}
}
