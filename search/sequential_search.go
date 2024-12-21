package search

import (
	"errors"
	"iter"
)

// SequentialSearchST represents an unordered symbol table of generic key-value pairs.
// It relies on the equals() function to test whether two keys are equal.
// This implementation uses a singly linked list and sequential search.
type SequentialSearchST[K, V any] struct {
	n      int                           // number of key-value pairs
	first  *sequentialSearchSTNode[K, V] // the linked list of key-value pairs
	equals func(a, b K) bool             // function to test whether two keys are equal
}

// sequentialSearchSTNode a helper linked list.
type sequentialSearchSTNode[K, V any] struct {
	key  K
	val  V
	next *sequentialSearchSTNode[K, V]
}

// NewSequentialSearchST initializes an empty symbol table.
// It gets a function as a parameter to test whether two keys are equal.
// The complexity is O(1).
func NewSequentialSearchST[K, V any](equals func(a, b K) bool) *SequentialSearchST[K, V] {
	return &SequentialSearchST[K, V]{
		n:      0,
		equals: equals,
	}
}

// Size returns the number of key-value pairs.
// The complexity is O(1).
func (s *SequentialSearchST[K, V]) Size() int {
	return s.n
}

// IsEmpty returns true if this symbol table is empty.
// The complexity is O(1).
func (s *SequentialSearchST[K, V]) IsEmpty() bool {
	return s.n == 0
}

// Contains returns true if this symbol table contains the specified key.
// The complexity is O(N) where N is the number of key-value pairs.
func (s *SequentialSearchST[K, V]) Contains(key K) bool {
	if _, err := s.Get(key); errors.Is(err, ErrAbsentKey) {
		return false
	}
	return true
}

// Get returns the value associated with the given key, ErrAbsentKey error if key is absent.
// The complexity is O(N) where N is the number of key-value pairs.
func (s *SequentialSearchST[K, V]) Get(key K) (V, error) {
	for x := s.first; x != nil; x = x.next {
		if s.equals(key, x.key) {
			return x.val, nil
		}
	}
	var value V
	return value, ErrAbsentKey
}

// Put Inserts the specified key-value pair, overwriting the old value with the new value if the symbol table
// already contains the specified key.
// The complexity is O(N) where N is the number of key-value pairs.
func (s *SequentialSearchST[K, V]) Put(key K, val V) {
	for x := s.first; x != nil; x = x.next {
		if s.equals(key, x.key) {
			x.val = val
			return
		}
	}
	s.first = &sequentialSearchSTNode[K, V]{
		key:  key,
		val:  val,
		next: s.first,
	}
	s.n++
}

// Delete removes the specified key and its associated value, ErrAbsentKey if key is absent.
// The complexity is O(N) where N is the number of key-value pairs.
func (s *SequentialSearchST[K, V]) Delete(key K) error {
	for p := &s.first; *p != nil; p = &(*p).next {
		if s.equals((*p).key, key) {
			*p = (*p).next
			s.n--
			return nil
		}
	}
	return ErrAbsentKey
}

// Iterator returns an iterator that iterates over all key-value pairs.
func (s *SequentialSearchST[K, V]) Iterator() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for node := s.first; node != nil; node = node.next {
			if !yield(node.key, node.val) {
				return
			}
		}
	}
}
