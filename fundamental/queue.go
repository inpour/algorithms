package fundamental

import (
	"errors"
	"sync"
)

// The queue type represents a first-in-first-out (FIFO) queue of generic items implemented
// using a linked list.
// It supports the usual Enqueue and Dequeue operations, along with methods for peeking at the
// first item, getting the size of the queue and testing if the queue is empty.
// The Enqueue, Dequeue, Peek, Size, and IsEmpty operations all take constant time in the worst case.
type queue struct {
	lock  *sync.Mutex // protect race condition
	first *queueNode  // beginning of queue
	last  *queueNode  // end of queue
	size  int         // number of items in  queue
}

// queueNode helper linked list.
type queueNode struct {
	item any
	next *queueNode
}

// NewQueue initializes an empty queue.
// The complexity is O(1).
func NewQueue() *queue {
	return &queue{
		lock: &sync.Mutex{},
		size: 0,
	}
}

var ErrEmptyQueue = errors.New("queue is empty")

// IsEmpty returns true if this queue is empty.
// The complexity is O(1).
func (queue *queue) IsEmpty() bool {
	return queue.first == nil
}

// Size returns the number of items in this queue.
// The complexity is O(1).
func (queue *queue) Size() int {
	return queue.size
}

// Enqueue adds the item to this queue.
// The complexity is O(1).
func (queue *queue) Enqueue(item any) {
	queue.lock.Lock()
	defer queue.lock.Unlock()

	oldLast := queue.last
	queue.last = &queueNode{item: item}
	queue.last.next = nil
	if queue.IsEmpty() {
		queue.first = queue.last
	} else {
		oldLast.next = queue.last
	}
	queue.size++
}

// Dequeue removes and returns the item least recently added to queue,
// returns ErrEmptyQueue if queue is empty.
// The complexity is O(1).
func (queue *queue) Dequeue() (any, error) {
	queue.lock.Lock()
	defer queue.lock.Unlock()

	if queue.IsEmpty() {
		return nil, ErrEmptyQueue
	}

	item := queue.first.item
	queue.first = queue.first.next
	if queue.IsEmpty() {
		queue.last = nil
	}
	queue.size--

	return item, nil
}

// Peek returns (but does not remove) the item least recently added to queue,
// returns ErrEmptyQueue if queue is empty.
// The complexity is O(1).
func (queue *queue) Peek() (any, error) {
	if queue.IsEmpty() {
		return nil, ErrEmptyQueue
	}

	return queue.first.item, nil
}
