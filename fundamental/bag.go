package fundamental

import (
	"iter"
	"sync"
)

// The bag type represents a bag (or multiset) of generic items. This implementation
// uses a linked list.
// It supports insertion and iterating over the items in arbitrary order.
// The Add, IsEmpty, and Size operations take constant time. Iteration takes time
// proportional to the number of items.
type bag[T any] struct {
	lock  *sync.Mutex // protect race condition
	first *bagNode[T] // beginning of bag
	size  int         // number of items in bag
}

// bagNode helper linked list.
type bagNode[T any] struct {
	item T
	next *bagNode[T]
}

// NewBag initializes an empty bag.
// The complexity is O(1).
func NewBag[T any]() *bag[T] {
	return &bag[T]{
		lock: &sync.Mutex{},
		size: 0,
	}
}

// IsEmpty returns true if this bag is empty.
// The complexity is O(1).
func (bag *bag[T]) IsEmpty() bool {
	return bag.first == nil
}

// Size returns the number of items in this bag.
// The complexity is O(1).
func (bag *bag[T]) Size() int {
	return bag.size
}

// Add the item to bag.
// The complexity is O(1).
func (bag *bag[T]) Add(item T) {
	bag.lock.Lock()
	defer bag.lock.Unlock()

	oldFirst := bag.first
	bag.first = &bagNode[T]{item: item, next: oldFirst}
	bag.size++
}

// Iterator returns an iterator that iterates over the items in the bag.
func (bag *bag[T]) Iterator() iter.Seq[T] {
	return func(yield func(T) bool) {
		if bag.first == nil {
			return
		}
		node := bag.first
		yield(node.item)
		for node.next != nil {
			node = node.next
			yield(node.item)
		}
	}
}
