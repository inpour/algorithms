package fundamental

import (
	"errors"
	"iter"
	"sync"
)

// Queue represents a first-in-first-out (FIFO) Queue of generic items. This implementation uses a singly linked list.
// It supports the usual Enqueue and Dequeue operations, along with methods for peeking at the
// first item, getting the size of the Queue and testing if the Queue is empty.
// Enqueue, Dequeue, Peek, Size, and IsEmpty perform O(1) work, excluding time waiting for the mutex.
// A Queue must not be copied after first use.
type Queue[T any] struct {
	lock  sync.Mutex    // protect race condition
	first *queueNode[T] // beginning of Queue
	last  *queueNode[T] // end of Queue
	size  int           // number of items in  Queue
}

// queueNode a helper linked list.
type queueNode[T any] struct {
	item T
	next *queueNode[T]
}

// NewQueue initializes an empty Queue.
// The complexity is O(1).
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		lock: sync.Mutex{},
		size: 0,
	}
}

var ErrEmptyQueue = errors.New("queue is empty")

// IsEmpty returns true if this Queue is empty.
// The complexity is O(1).
func (queue *Queue[T]) IsEmpty() bool {
	queue.lock.Lock()
	defer queue.lock.Unlock()

	return queue.first == nil
}

// isEmptyLocked reports whether the queue is empty.
// The caller must hold queue.lock.
func (queue *Queue[T]) isEmptyLocked() bool {
	return queue.first == nil
}

// Size returns the number of items in this Queue.
// The complexity is O(1).
func (queue *Queue[T]) Size() int {
	queue.lock.Lock()
	defer queue.lock.Unlock()

	return queue.size
}

// Enqueue adds the item to this Queue.
// The complexity is O(1).
func (queue *Queue[T]) Enqueue(item T) {
	queue.lock.Lock()
	defer queue.lock.Unlock()

	oldLast := queue.last
	queue.last = &queueNode[T]{item: item, next: nil}
	if queue.isEmptyLocked() {
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
	if queue.isEmptyLocked() {
		return item, ErrEmptyQueue
	}
	item = queue.first.item
	queue.first = queue.first.next
	if queue.isEmptyLocked() {
		queue.last = nil
	}
	queue.size--

	return item, nil
}

// Peek returns (but does not remove) the item least recently added to Queue,
// returns ErrEmptyQueue if Queue is empty.
// The complexity is O(1).
func (queue *Queue[T]) Peek() (T, error) {
	queue.lock.Lock()
	defer queue.lock.Unlock()

	var item T
	if queue.isEmptyLocked() {
		return item, ErrEmptyQueue
	}

	return queue.first.item, nil
}

// Iterator returns an iterator over the queue's items in FIFO order
// without taking a snapshot or acquiring the queue's mutex.
//
// The caller must ensure that no Enqueue or Dequeue operation overlaps
// traversal, including operations called from the loop body.
// Any preceding queue modifications must be synchronized with traversal.
//
// A full traversal takes O(n) time and uses O(1) additional space.
func (queue *Queue[T]) Iterator() iter.Seq[T] {
	return func(yield func(T) bool) {
		for node := queue.first; node != nil; node = node.next {
			if !yield(node.item) {
				return
			}
		}
	}
}

// IteratorSnapshot returns an iterator over a FIFO-ordered snapshot captured when traversal begins.
// Capturing the snapshot takes O(n) time and O(n) additional space.
// Item values are copied; referenced objects are not deep-copied.
func (queue *Queue[T]) IteratorSnapshot() iter.Seq[T] {
	return func(yield func(T) bool) {
		items := func() []T {
			queue.lock.Lock()
			defer queue.lock.Unlock()

			snapshot := make([]T, 0, queue.size)
			for node := queue.first; node != nil; node = node.next {
				snapshot = append(snapshot, node.item)
			}

			return snapshot
		}()

		for _, item := range items {
			if !yield(item) {
				return
			}
		}
	}
}
