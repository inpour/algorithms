package fundamental

import (
	"errors"
	"iter"
	"sync"
)

// The Stack type represents a last-in-first-out (LIFO) Stack of generic items implemented
// using a linked list.
// It supports the usual Push and Pop operations, along with methods for Peeking at the
// top item, getting the size of the Stack and testing if the Stack is empty.
// The Push, Pop, Peek, Size, and IsEmpty operations all take constant time in the worst case.
type Stack[T any] struct {
	lock  *sync.Mutex   // protect race condition
	first *stackNode[T] // top of Stack
	size  int           // number of items in Stack
}

// stackNode helper linked list.
type stackNode[T any] struct {
	item T
	next *stackNode[T]
}

// NewStack initializes an empty Stack.
// The complexity is O(1).
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		lock: &sync.Mutex{},
		size: 0,
	}
}

var ErrEmptyStack = errors.New("stack is empty")

// IsEmpty returns true if this Stack is empty.
// The complexity is O(1).
func (stack *Stack[T]) IsEmpty() bool {
	return stack.first == nil
}

// Size returns the number of items in this Stack.
// The complexity is O(1).
func (stack *Stack[T]) Size() int {
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
	if stack.IsEmpty() {
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
	var item T
	if stack.IsEmpty() {
		return item, ErrEmptyStack
	}

	return stack.first.item, nil
}

// Iterator returns an iterator that iterates over the items in the Stack.
func (stack *Stack[T]) Iterator() iter.Seq[T] {
	return func(yield func(T) bool) {
		for node := stack.first; node != nil; node = node.next {
			if !yield(node.item) {
				return
			}
		}
	}
}
