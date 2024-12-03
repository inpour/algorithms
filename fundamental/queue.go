package fundamental

import (
	"errors"
	"sync"
)

// The Queue type represents a first-in-first-out (FIFO) Queue of generic items implemented
// using a linked list.
// It supports the usual Enqueue and Dequeue operations, along with methods for peeking at the
// first item, getting the size of the Queue and testing if the Queue is empty.
// The Enqueue, Dequeue, Peek, Size, and IsEmpty operations all take constant time in the worst case.
type Queue struct {
	lock  *sync.Mutex // protect race condition
	first *queueNode  // beginning of Queue
	last  *queueNode  // end of Queue
	size  int         // number of items in  Queue
}

// queueNode helper linked list.
type queueNode struct {
	item any
	next *queueNode
}

// NewQueue initializes an empty Queue.
// The complexity is O(1).
func NewQueue() *Queue {
	return &Queue{
		lock: &sync.Mutex{},
		size: 0,
	}
}

var ErrEmptyQueue = errors.New("queue is empty")

// IsEmpty returns true if this Queue is empty.
// The complexity is O(1).
func (queue *Queue) IsEmpty() bool {
	return queue.first == nil
}

// Size returns the number of items in this Queue.
// The complexity is O(1).
func (queue *Queue) Size() int {
	return queue.size
}

// Enqueue adds the item to this Queue.
// The complexity is O(1).
func (queue *Queue) Enqueue(item any) {
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

// Dequeue removes and returns the item least recently added to Queue,
// returns ErrEmptyQueue if Queue is empty.
// The complexity is O(1).
func (queue *Queue) Dequeue() (any, error) {
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

// Peek returns (but does not remove) the item least recently added to Queue,
// returns ErrEmptyQueue if Queue is empty.
// The complexity is O(1).
func (queue *Queue) Peek() (any, error) {
	if queue.IsEmpty() {
		return nil, ErrEmptyQueue
	}

	return queue.first.item, nil
}
