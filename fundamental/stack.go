package fundamental

import (
	"errors"
	"sync"
)

// The stack type represents a last-in-first-out (LIFO) stack of generic items implemented
// using a linked list.
// It supports the usual Push and Pop operations, along with methods for Peeking at the
// top item, getting the size of the stack and testing if the stack is empty.
// The Push, Pop, Peek, Size, and IsEmpty operations all take constant time in the worst case.
type stack struct {
	lock  *sync.Mutex // protect race condition
	first *stackNode  // top of stack
	size  int         // number of items in stack
}

// stackNode helper linked list.
type stackNode struct {
	item any
	next *stackNode
}

// NewStack initializes an empty stack.
// The complexity is O(1).
func NewStack() *stack {
	return &stack{
		lock: &sync.Mutex{},
		size: 0,
	}
}

var ErrEmptyStack = errors.New("stack is empty")

// IsEmpty returns true if this stack is empty.
// The complexity is O(1).
func (stack *stack) IsEmpty() bool {
	return stack.first == nil
}

// Size returns the number of items in this stack.
// The complexity is O(1).
func (stack *stack) Size() int {
	return stack.size
}

// Push adds the item to stack.
// The complexity is O(1).
func (stack *stack) Push(item any) {
	stack.lock.Lock()
	defer stack.lock.Unlock()

	oldFirst := stack.first
	stack.first = &stackNode{item: item}
	stack.first.next = oldFirst
	stack.size++
}

// Pop removes and returns the item most recently added to stack,
// returns ErrEmptyStack if stack is empty.
// The complexity is O(1).
func (stack *stack) Pop() (any, error) {
	stack.lock.Lock()
	defer stack.lock.Unlock()

	if stack.IsEmpty() {
		return nil, ErrEmptyStack
	}
	item := stack.first.item
	stack.first = stack.first.next
	stack.size--

	return item, nil
}

// Peek returns (but does not remove) the item most recently added to stack,
// returns ErrEmptyStack if stack is empty.
// The complexity is O(1).
func (stack *stack) Peek() (any, error) {
	if stack.IsEmpty() {
		return nil, ErrEmptyStack
	}

	return stack.first.item, nil
}
