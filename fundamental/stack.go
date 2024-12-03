package fundamental

import (
	"errors"
	"sync"
)

// The Stack type represents a last-in-first-out (LIFO) Stack of generic items implemented
// using a linked list.
// It supports the usual Push and Pop operations, along with methods for Peeking at the
// top item, getting the size of the Stack and testing if the Stack is empty.
// The Push, Pop, Peek, Size, and IsEmpty operations all take constant time in the worst case.
type Stack struct {
	lock  *sync.Mutex // protect race condition
	first *stackNode  // top of Stack
	size  int         // number of items in Stack
}

// stackNode helper linked list.
type stackNode struct {
	item any
	next *stackNode
}

// NewStack initializes an empty Stack.
// The complexity is O(1).
func NewStack() *Stack {
	return &Stack{
		lock: &sync.Mutex{},
		size: 0,
	}
}

var ErrEmptyStack = errors.New("stack is empty")

// IsEmpty returns true if this Stack is empty.
// The complexity is O(1).
func (stack *Stack) IsEmpty() bool {
	return stack.first == nil
}

// Size returns the number of items in this Stack.
// The complexity is O(1).
func (stack *Stack) Size() int {
	return stack.size
}

// Push adds the item to Stack.
// The complexity is O(1).
func (stack *Stack) Push(item any) {
	stack.lock.Lock()
	defer stack.lock.Unlock()

	oldFirst := stack.first
	stack.first = &stackNode{item: item}
	stack.first.next = oldFirst
	stack.size++
}

// Pop removes and returns the item most recently added to Stack,
// returns ErrEmptyStack if Stack is empty.
// The complexity is O(1).
func (stack *Stack) Pop() (any, error) {
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

// Peek returns (but does not remove) the item most recently added to Stack,
// returns ErrEmptyStack if Stack is empty.
// The complexity is O(1).
func (stack *Stack) Peek() (any, error) {
	if stack.IsEmpty() {
		return nil, ErrEmptyStack
	}

	return stack.first.item, nil
}
