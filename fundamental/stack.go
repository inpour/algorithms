package fundamental

import (
	"errors"
	"iter"
	"sync"
)

// Stack represents a last-in-first-out (LIFO) Stack of generic items. This implementation uses a singly linked list.
// It supports the usual Push and Pop operations, along with methods for Peeking at the
// top item, getting the size of the Stack and testing if the Stack is empty.
// Push, Pop, Peek, Size, and IsEmpty perform O(1) work, excluding time waiting for the mutex.
// A Stack must not be copied after first use.
type Stack[T any] struct {
	lock  sync.Mutex    // protect race condition
	first *stackNode[T] // top of Stack
	size  int           // number of items in Stack
}

// stackNode a helper linked list.
type stackNode[T any] struct {
	item T
	next *stackNode[T]
}

// NewStack initializes an empty Stack.
// The complexity is O(1).
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		lock: sync.Mutex{},
		size: 0,
	}
}

var ErrEmptyStack = errors.New("stack is empty")

// IsEmpty returns true if this Stack is empty.
// The complexity is O(1).
func (stack *Stack[T]) IsEmpty() bool {
	stack.lock.Lock()
	defer stack.lock.Unlock()

	return stack.first == nil
}

// isEmptyLocked reports whether the stack is empty.
// The caller must hold stack.lock.
func (stack *Stack[T]) isEmptyLocked() bool {
	return stack.first == nil
}

// Size returns the number of items in this Stack.
// The complexity is O(1).
func (stack *Stack[T]) Size() int {
	stack.lock.Lock()
	defer stack.lock.Unlock()

	return stack.size
}

// Push adds the item to Stack.
// The complexity is O(1).
func (stack *Stack[T]) Push(item T) {
	stack.lock.Lock()
	defer stack.lock.Unlock()

	oldFirst := stack.first
	stack.first = &stackNode[T]{item: item, next: oldFirst}
	stack.size++
}

// Pop removes and returns the item most recently added to Stack,
// returns ErrEmptyStack if Stack is empty.
// The complexity is O(1).
func (stack *Stack[T]) Pop() (T, error) {
	stack.lock.Lock()
	defer stack.lock.Unlock()

	var item T
	if stack.isEmptyLocked() {
		return item, ErrEmptyStack
	}
	item = stack.first.item
	stack.first = stack.first.next
	stack.size--

	return item, nil
}

// Peek returns (but does not remove) the item most recently added to Stack,
// returns ErrEmptyStack if Stack is empty.
// The complexity is O(1).
func (stack *Stack[T]) Peek() (T, error) {
	stack.lock.Lock()
	defer stack.lock.Unlock()

	var item T
	if stack.isEmptyLocked() {
		return item, ErrEmptyStack
	}

	return stack.first.item, nil
}

// Iterator returns an iterator over the stack's items in LIFO order.
// Each traversal captures the top node when traversal begins.
// Later pushes are excluded, and later pops do not remove items
// from the captured traversal.
//
// Capturing the top node takes O(1) work, excluding mutex wait time.
// A full traversal takes O(n) time and uses O(1) additional space.
// The stack's mutex is not held while yielding items.
//
// Item values are yielded without deep-copying referenced objects.
func (stack *Stack[T]) Iterator() iter.Seq[T] {
	return func(yield func(T) bool) {
		stack.lock.Lock()
		first := stack.first
		stack.lock.Unlock()

		for node := first; node != nil; node = node.next {
			if !yield(node.item) {
				return
			}
		}
	}
}
