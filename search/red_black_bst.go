package search

import (
	"iter"
)

// colors of parent link
const (
	red   = true
	black = false
)

// RedBlackBST represents an ordered symbol table of generic key-value pairs.
// It relies on the compare() function to compare two keys:
//
//	if a == b then compare(a, b) returns 0
//	if a > b then compare(a, b) returns 1
//	if a < b then compare(a, b) returns -1
//
// This implementation uses a left-leaning red-black binary search tree.
type RedBlackBST[K, V any] struct {
	root    *redBlackBSTNode[K, V] // root of BST
	compare func(a, b K) int       // function to compare two keys
}

// redBlackBSTNode a helper linked list.
type redBlackBSTNode[K, V any] struct {
	key         K                      // sorted by key
	val         V                      // associated data
	left, right *redBlackBSTNode[K, V] // left and right subtrees
	color       bool                   // color of parent link (red or black)
	size        int                    // number of nodes in subtree
}

// NewRedBlackBST initializes an empty symbol table.
// It gets a function as a parameter to compare two keys.
// The complexity is O(1).
func NewRedBlackBST[K, V any](compare func(a, b K) int) *RedBlackBST[K, V] {
	return &RedBlackBST[K, V]{
		compare: compare,
	}
}

// IsEmpty returns true if this symbol table is empty.
// The complexity is O(1).
func (b *RedBlackBST[K, V]) IsEmpty() bool {
	return b.Size() == 0
}

// Size returns the number of key-value pairs.
// The complexity is O(1).
func (b *RedBlackBST[K, V]) Size() int {
	return b.size(b.root)
}

// size return number of key-value pairs in RedBlackBST rooted at node.
// The complexity is O(1).
func (b *RedBlackBST[K, V]) size(node *redBlackBSTNode[K, V]) int {
	if node == nil {
		return 0
	}
	return node.size
}

// Contains returns true if this symbol table contain the given key.
// The complexity is O(log(N)) where N is the number of key-value pairs.
func (b *RedBlackBST[K, V]) Contains(key K) bool {
	_, err := b.Get(key)
	return err == nil
}

// Get returns the value associated with the given key, ErrAbsentKey if key is absent.
// The complexity is O(log(N)) where N is the number of key-value pairs.
func (b *RedBlackBST[K, V]) Get(key K) (V, error) {
	return b.get(b.root, key)
}

func (b *RedBlackBST[K, V]) get(node *redBlackBSTNode[K, V], key K) (V, error) {
	for node != nil {
		cmp := b.compare(key, node.key)
		if cmp < 0 {
			node = node.left
		} else if cmp > 0 {
			node = node.right
		} else {
			return node.val, nil
		}
	}
	var value V
	return value, ErrAbsentKey
}

// Put inserts the specified key-value pair into the symbol table, overwriting the old value with the
// new value if the symbol table already contains the specified key.
// The complexity is O(log(N)) where N is the number of key-value pairs.
func (b *RedBlackBST[K, V]) Put(key K, val V) {
	b.root = b.put(b.root, key, val)
	b.root.color = black
}

func (b *RedBlackBST[K, V]) put(node *redBlackBSTNode[K, V], key K, val V) *redBlackBSTNode[K, V] {
	if node == nil {
		return &redBlackBSTNode[K, V]{
			key:   key,
			val:   val,
			color: red,
			size:  1,
		}
	}

	cmp := b.compare(key, node.key)
	if cmp < 0 {
		node.left = b.put(node.left, key, val)
	} else if cmp > 0 {
		node.right = b.put(node.right, key, val)
	} else {
		node.val = val
	}

	// fix-up any right-leaning links
	return b.balance(node)
}

// DelMin removes the smallest key and associated value, ErrEmptySymbolTable if the symbol table is empty.
// The complexity is O(log(N)) where N is the number of key-value pairs.
func (b *RedBlackBST[K, V]) DelMin() error {
	if b.IsEmpty() {
		return ErrEmptySymbolTable
	}

	// if both children of root are black, set root to red
	if !b.isRed(b.root.left) && !b.isRed(b.root.right) {
		b.root.color = red
	}

	b.root = b.delMin(b.root)
	if !b.IsEmpty() {
		b.root.color = black
	}
	return nil
}

func (b *RedBlackBST[K, V]) delMin(node *redBlackBSTNode[K, V]) *redBlackBSTNode[K, V] {
	if node.left == nil {
		return nil
	}

	if !b.isRed(node.left) && !b.isRed(node.left.left) {
		node = b.moveRedLeft(node)
	}

	node.left = b.delMin(node.left)
	return b.balance(node)
}

// DelMax removes the largest key and associated value, ErrEmptySymbolTable if the symbol table is empty.
// The complexity is O(log(N)) where N is the number of key-value pairs.
func (b *RedBlackBST[K, V]) DelMax() error {
	if b.IsEmpty() {
		return ErrEmptySymbolTable
	}

	// if both children of root are black, set root to red
	if !b.isRed(b.root.left) && !b.isRed(b.root.right) {
		b.root.color = red
	}

	b.root = b.delMax(b.root)
	if !b.IsEmpty() {
		b.root.color = black
	}
	return nil
}

func (b *RedBlackBST[K, V]) delMax(node *redBlackBSTNode[K, V]) *redBlackBSTNode[K, V] {
	if b.isRed(node.left) {
		node = b.rotateRight(node)
	}

	if node.right == nil {
		return nil
	}

	if !b.isRed(node.right) && !b.isRed(node.right.left) {
		node = b.moveRedRight(node)
	}

	node.right = b.delMax(node.right)
	return b.balance(node)
}

// Min returns the smallest key, ErrEmptySymbolTable if the symbol table is empty.
// The complexity is O(log(N)) where N is the number of key-value pairs.
func (b *RedBlackBST[K, V]) Min() (K, error) {
	if b.IsEmpty() {
		var key K
		return key, ErrEmptySymbolTable
	}
	return b.min(b.root).key, nil
}

func (b *RedBlackBST[K, V]) min(node *redBlackBSTNode[K, V]) *redBlackBSTNode[K, V] {
	if node.left == nil {
		return node
	}
	return b.min(node.left)
}

// Max returns the largest key, ErrEmptySymbolTable if the symbol table is empty.
// The complexity is O(log(N)) where N is the number of key-value pairs.
func (b *RedBlackBST[K, V]) Max() (K, error) {
	if b.IsEmpty() {
		var key K
		return key, ErrEmptySymbolTable
	}
	return b.max(b.root).key, nil
}

func (b *RedBlackBST[K, V]) max(node *redBlackBSTNode[K, V]) *redBlackBSTNode[K, V] {
	if node.right == nil {
		return node
	}
	return b.max(node.right)
}

// Delete removes the specified key and associated value, ErrAbsentKey if key is absent.
// The complexity is O(log(N)) where N is the number of key-value pairs.
func (b *RedBlackBST[K, V]) Delete(key K) error {
	if !b.Contains(key) {
		return ErrAbsentKey
	}

	// if both children of root are black, set root to red
	if !b.isRed(b.root.left) && !b.isRed(b.root.right) {
		b.root.color = red
	}

	b.root = b.delete(b.root, key)
	return nil
}

func (b *RedBlackBST[K, V]) delete(node *redBlackBSTNode[K, V], key K) *redBlackBSTNode[K, V] {
	if b.compare(key, node.key) < 0 {
		if !b.isRed(node.left) && !b.isRed(node.left.left) {
			node = b.moveRedLeft(node)
		}
		node.left = b.delete(node.left, key)
	} else {
		if b.isRed(node.left) {
			node = b.rotateRight(node)
		}
		if b.compare(key, node.key) == 0 && node.right == nil {
			return nil
		}
		if !b.isRed(node.right) && !b.isRed(node.right.left) {
			node = b.moveRedRight(node)
		}
		if b.compare(key, node.key) == 0 {
			x := b.min(node.right)
			node.key = x.key
			node.val = x.val
			node.right = b.delMin(node.right)
		} else {
			node.right = b.delete(node.right, key)
		}
	}
	return b.balance(node)
}

// Floor returns the largest key less than or equal to key, ErrTooSmallFloorKey if key to floor is too small.
// The complexity is O(log(N)) where N is the number of key-value pairs.
func (b *RedBlackBST[K, V]) Floor(key K) (K, error) {
	node, err := b.floor(b.root, key)
	if err != nil {
		return key, err
	}
	return node.key, nil
}

func (b *RedBlackBST[K, V]) floor(node *redBlackBSTNode[K, V], key K) (*redBlackBSTNode[K, V], error) {
	if node == nil {
		return nil, ErrTooSmallFloorKey
	}
	cmp := b.compare(key, node.key)
	if cmp == 0 {
		return node, nil
	}
	if cmp < 0 {
		return b.floor(node.left, key)
	}
	if tmpNode, err := b.floor(node.right, key); err == nil {
		return tmpNode, nil
	}
	return node, nil
}

// Ceiling returns the smallest key greater than or equal to key, ErrTooLargeCeilingKey if key to ceiling is too large.
// The complexity is O(log(N)) where N is the number of key-value pairs.
func (b *RedBlackBST[K, V]) Ceiling(key K) (K, error) {
	node, err := b.ceiling(b.root, key)
	if err != nil {
		return key, err
	}
	return node.key, nil
}

func (b *RedBlackBST[K, V]) ceiling(node *redBlackBSTNode[K, V], key K) (*redBlackBSTNode[K, V], error) {
	if node == nil {
		return nil, ErrTooLargeCeilingKey
	}
	cmp := b.compare(key, node.key)
	if cmp == 0 {
		return node, nil
	}
	if cmp > 0 {
		return b.ceiling(node.right, key)
	}
	if tmpNode, err := b.ceiling(node.left, key); err == nil {
		return tmpNode, nil
	}
	return node, nil
}

// Select return the kth smallest key (key of rank k), ErrInvalidRank if rank is out of range.
// The complexity is O(log(N)) where N is the number of key-value pairs.
func (b *RedBlackBST[K, V]) Select(k int) (K, error) {
	var key K
	if k < 0 || k >= b.Size() {
		return key, ErrInvalidRank
	}
	return b.selectRecursive(b.root, k)
}

func (b *RedBlackBST[K, V]) selectRecursive(node *redBlackBSTNode[K, V], k int) (K, error) {
	leftSize := b.size(node.left)
	if leftSize > k {
		return b.selectRecursive(node.left, k)
	} else if leftSize < k {
		return b.selectRecursive(node.right, k-leftSize-1)
	} else {
		return node.key, nil
	}
}

// Rank returns the number of keys strictly less than key, ErrAbsentKey if key is absent.
// The complexity is O(log(N)) where N is the number of key-value pairs.
func (b *RedBlackBST[K, V]) Rank(key K) (int, error) {
	return b.rank(b.root, key)
}

func (b *RedBlackBST[K, V]) rank(node *redBlackBSTNode[K, V], key K) (int, error) {
	if node == nil {
		return 0, ErrAbsentKey
	}
	cmp := b.compare(key, node.key)
	if cmp < 0 {
		return b.rank(node.left, key)
	} else if cmp > 0 {
		rightRank, err := b.rank(node.right, key)
		return 1 + b.size(node.left) + rightRank, err
	} else {
		return b.size(node.left), nil
	}
}

// Iterator returns an iterator that iterates over all key-value pairs in sorted order.
// It takes O(log(N)) time to prepare iterator where N is the number of key-value pairs.
func (b *RedBlackBST[K, V]) Iterator() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {

		lo, err := b.Min()
		if err != nil {
			return
		}
		hi, _ := b.Max()
		b.iterator(yield, b.root, lo, hi)
	}

}

// RangeIterator returns an iterator that iterates over key-value pairs where keys in [lo:hi] range, in sorted order.
// It takes O(log(N)) time to prepare iterator where N is the number of key-value pairs.
func (b *RedBlackBST[K, V]) RangeIterator(lo, hi K) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if b.compare(lo, hi) > 0 {
			return
		}
		b.iterator(yield, b.root, lo, hi)
	}
}

func (b *RedBlackBST[K, V]) iterator(yield func(K, V) bool, node *redBlackBSTNode[K, V], lo, hi K) {
	if node == nil {
		return
	}
	cmpLo := b.compare(lo, node.key)
	cmpHi := b.compare(hi, node.key)
	if cmpLo < 0 {
		b.iterator(yield, node.left, lo, hi)
	}
	if cmpLo <= 0 && cmpHi >= 0 {
		if !yield(node.key, node.val) {
			return
		}
	}
	if cmpHi > 0 {
		b.iterator(yield, node.right, lo, hi)
	}
}

// RangeSize returns the number of keys in [lo:hi] range.
// The complexity is O(log(N)) where N is the number of key-value pairs.
func (b *RedBlackBST[K, V]) RangeSize(lo, hi K) int {
	if b.compare(lo, hi) > 0 {
		return 0
	}
	hiRank, err := b.Rank(hi)
	loRank, _ := b.Rank(lo)
	size := hiRank - loRank
	// symbol table contains hi key
	if err == nil {
		size++
	}
	return size
}

// isRed returns true if node is red, false if node is nil.
func (b *RedBlackBST[K, V]) isRed(node *redBlackBSTNode[K, V]) bool {
	if node == nil {
		return false
	}
	return node.color == red
}

// rotateRight makes a left-leaning link lean to the right.
func (b *RedBlackBST[K, V]) rotateRight(node *redBlackBSTNode[K, V]) *redBlackBSTNode[K, V] {
	x := node.left
	node.left = x.right
	x.right = node
	x.color = node.color
	node.color = red
	x.size = node.size
	node.size = 1 + b.size(node.left) + b.size(node.right)
	return x
}

// rotateLeft makes a right-leaning link lean to the left.
func (b *RedBlackBST[K, V]) rotateLeft(node *redBlackBSTNode[K, V]) *redBlackBSTNode[K, V] {
	x := node.right
	node.right = x.left
	x.left = node
	x.color = node.color
	node.color = red
	x.size = node.size
	node.size = 1 + b.size(node.left) + b.size(node.right)
	return x
}

// flipColors flips the colors of a node and its two children.
func (b *RedBlackBST[K, V]) flipColors(node *redBlackBSTNode[K, V]) {
	node.color = !node.color
	node.left.color = !node.left.color
	node.right.color = !node.right.color
}

// moveRedLeft makes node.left or one of its children red,
// assuming that node is red and both node.left and node.left.left are black.
func (b *RedBlackBST[K, V]) moveRedLeft(node *redBlackBSTNode[K, V]) *redBlackBSTNode[K, V] {
	b.flipColors(node)
	if b.isRed(node.right.left) {
		node.right = b.rotateRight(node.right)
		node = b.rotateLeft(node)
		b.flipColors(node)
	}
	return node
}

// moveRedRight makes node.right or one of its children red,
// assuming that node is red and both node.right and node.right.left are black.
func (b *RedBlackBST[K, V]) moveRedRight(node *redBlackBSTNode[K, V]) *redBlackBSTNode[K, V] {
	b.flipColors(node)
	if b.isRed(node.left.left) {
		node = b.rotateRight(node)
		b.flipColors(node)
	}
	return node
}

// balance restores red-black tree invariant.
func (b *RedBlackBST[K, V]) balance(node *redBlackBSTNode[K, V]) *redBlackBSTNode[K, V] {
	if b.isRed(node.right) && !b.isRed(node.left) {
		node = b.rotateLeft(node)
	}
	if b.isRed(node.left) && b.isRed(node.left.left) {
		node = b.rotateRight(node)
	}
	if b.isRed(node.left) && b.isRed(node.right) {
		b.flipColors(node)
	}
	node.size = 1 + b.size(node.left) + b.size(node.right)
	return node
}
