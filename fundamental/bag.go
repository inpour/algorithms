package fundamental

import (
	"iter"
	"sync"
)

// Bag represents a Bag (or multiset) of generic items. This implementation uses a singly linked list.
// It supports insertion and iterating over the items in arbitrary order.
// The Add, IsEmpty, and Size operations take constant time. Iteration takes time
// proportional to the number of items.
type Bag[T any] struct {
	lock  *sync.Mutex // protect race condition
	first *bagNode[T] // beginning of Bag
	size  int         // number of items in Bag
}

// bagNode a helper linked list.
type bagNode[T any] struct {
	item T
	next *bagNode[T]
}

// NewBag initializes an empty Bag.
// The complexity is O(1).
func NewBag[T any]() *Bag[T] {
	return &Bag[T]{
		lock: &sync.Mutex{},
		size: 0,
	}
}

// IsEmpty returns true if this Bag is empty.
// The complexity is O(1).
func (bag *Bag[T]) IsEmpty() bool {
	return bag.first == nil
}

// Size returns the number of items in this Bag.
// The complexity is O(1).
func (bag *Bag[T]) Size() int {
	return bag.size
}

// Add the item to Bag.
// The complexity is O(1).
func (bag *Bag[T]) Add(item T) {
	bag.lock.Lock()
	defer bag.lock.Unlock()

	oldFirst := bag.first
	bag.first = &bagNode[T]{item: item, next: oldFirst}
	bag.size++
}

// Iterator returns an iterator that iterates over the items in the Bag.
func (bag *Bag[T]) Iterator() iter.Seq[T] {
	return func(yield func(T) bool) {
		for node := bag.first; node != nil; node = node.next {
			if !yield(node.item) {
				return
			}
		}
	}
}
