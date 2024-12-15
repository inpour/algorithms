package fundamental

import (
	"errors"
	"iter"
	"sync"
)

// The Queue type represents a first-in-first-out (FIFO) Queue of generic items implemented
// using a linked list.
// It supports the usual Enqueue and Dequeue operations, along with methods for peeking at the
// first item, getting the size of the Queue and testing if the Queue is empty.
// The Enqueue, Dequeue, Peek, Size, and IsEmpty operations all take constant time in the worst case.
type Queue[T any] struct {
	lock  *sync.Mutex   // protect race condition
	first *queueNode[T] // beginning of Queue
	last  *queueNode[T] // end of Queue
	size  int           // number of items in  Queue
}

// queueNode helper linked list.
type queueNode[T any] struct {
	item T
	next *queueNode[T]
}

// NewQueue initializes an empty Queue.[T
// The complexity is O(1).
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		lock: &sync.Mutex{},
		size: 0,
	}
}

var ErrEmptyQueue = errors.New("queue is empty")

// IsEmpty returns true if this Queue is empty.
// The complexity is O(1).
func (queue *Queue[T]) IsEmpty() bool {
	return queue.first == nil
}

// Size returns the number of items in this Queue.
// The complexity is O(1).
func (queue *Queue[T]) Size() int {
	return queue.size
}

// Enqueue adds the item to this Queue.
// The complexity is O(1).
func (queue *Queue[T]) Enqueue(item T) {
	queue.lock.Lock()
	defer queue.lock.Unlock()

	oldLast := queue.last
	queue.last = &queueNode[T]{item: item, next: nil}
	if queue.IsEmpty() {
		queue.first = queue.last
	} else {
		oldLast.next = queue.last
	}
	queue.size++
}

// Dequeue removes and returns the item least recently added to Queue,
// returns ErrEmptyQueue if Queue is empty.
// The complexity is O(1).
func (queue *Queue[T]) Dequeue() (T, error) {
	queue.lock.Lock()
	defer queue.lock.Unlock()

	var item T
	if queue.IsEmpty() {
		return item, ErrEmptyQueue
	}
	item = queue.first.item
	queue.first = queue.first.next
	if queue.IsEmpty() {
		queue.last = nil
	}
	queue.size--

	return item, nil
}

// Peek returns (but does not remove) the item least recently added to Queue,
// returns ErrEmptyQueue if Queue is empty.
// The complexity is O(1).
func (queue *Queue[T]) Peek() (T, error) {
	var item T
	if queue.IsEmpty() {
		return item, ErrEmptyQueue
	}

	return queue.first.item, nil
}

// Iterator returns an iterator that iterates over the items in the Queue.
func (queue *Queue[T]) Iterator() iter.Seq[T] {
	return func(yield func(T) bool) {
		for node := queue.first; node != nil; node = node.next {
			if !yield(node.item) {
				return
			}
		}
	}
}
