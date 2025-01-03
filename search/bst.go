package search

import (
	"iter"
)

// BST represents an ordered symbol table of generic key-value pairs.
// It relies on the compare() function to compare two keys:
//
//	if a == b then compare(a, b) returns 0
//	if a > b then compare(a, b) returns 1
//	if a < b then compare(a, b) returns -1
//
// This implementation uses an unbalanced binary search tree.
type BST[K, V any] struct {
	root    *bstNode[K, V]   // root of BST
	compare func(a, b K) int // function to compare two keys
}

// bstNode a helper linked list.
type bstNode[K, V any] struct {
	key         K              // sorted by key
	val         V              // associated data
	left, right *bstNode[K, V] // left and right subtrees
	size        int            // number of nodes in subtree
}

// NewBST initializes an empty symbol table.
// It gets a function as a parameter to compare two keys.
// The complexity is O(1).
func NewBST[K, V any](compare func(a, b K) int) *BST[K, V] {
	return &BST[K, V]{
		compare: compare,
	}
}

// IsEmpty returns true if this symbol table is empty.
// The complexity is O(1).
func (b *BST[K, V]) IsEmpty() bool {
	return b.Size() == 0
}

// Size returns the number of key-value pairs.
// The complexity is O(1).
func (b *BST[K, V]) Size() int {
	return b.size(b.root)
}

// size return number of key-value pairs in BST rooted at node.
// The complexity is O(1).
func (b *BST[K, V]) size(node *bstNode[K, V]) int {
	if node == nil {
		return 0
	}
	return node.size
}

// Contains returns true if this symbol table contain the given key.
// The complexity is O(N) where N is the number of key-value pairs.
func (b *BST[K, V]) Contains(key K) bool {
	_, err := b.Get(key)
	return err == nil
}

// Get returns the value associated with the given key, ErrAbsentKey if key is absent.
// The complexity is O(N) where N is the number of key-value pairs.
func (b *BST[K, V]) Get(key K) (V, error) {
	return b.get(b.root, key)
}

func (b *BST[K, V]) get(node *bstNode[K, V], key K) (V, error) {
	if node == nil {
		var value V
		return value, ErrAbsentKey
	}

	cmp := b.compare(key, node.key)
	if cmp < 0 {
		return b.get(node.left, key)
	} else if cmp > 0 {
		return b.get(node.right, key)
	} else {
		return node.val, nil
	}
}

// Put inserts the specified key-value pair into the symbol table, overwriting the old value with the
// new value if the symbol table already contains the specified key.
// The complexity is O(N) where N is the number of key-value pairs.
func (b *BST[K, V]) Put(key K, val V) {
	b.root = b.put(b.root, key, val)
}

func (b *BST[K, V]) put(node *bstNode[K, V], key K, val V) *bstNode[K, V] {
	if node == nil {
		return &bstNode[K, V]{
			key:  key,
			val:  val,
			size: 1,
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
	node.size = 1 + b.size(node.left) + b.size(node.right)
	return node
}

// DelMin removes the smallest key and associated value, ErrEmptySymbolTable if the symbol table is empty.
// The complexity is O(N) where N is the number of key-value pairs.
func (b *BST[K, V]) DelMin() error {
	if b.IsEmpty() {
		return ErrEmptySymbolTable
	}
	b.root = b.delMin(b.root)
	return nil
}

func (b *BST[K, V]) delMin(node *bstNode[K, V]) *bstNode[K, V] {
	if node.left == nil {
		return node.right
	}
	node.left = b.delMin(node.left)
	node.size = 1 + b.size(node.left) + b.size(node.right)
	return node
}

// DelMax removes the largest key and associated value, ErrEmptySymbolTable if the symbol table is empty.
// The complexity is O(N) where N is the number of key-value pairs.
func (b *BST[K, V]) DelMax() error {
	if b.IsEmpty() {
		return ErrEmptySymbolTable
	}
	b.root = b.delMax(b.root)
	return nil
}

func (b *BST[K, V]) delMax(node *bstNode[K, V]) *bstNode[K, V] {
	if node.right == nil {
		return node.left
	}
	node.right = b.delMax(node.right)
	node.size = 1 + b.size(node.left) + b.size(node.right)
	return node
}

// Min returns the smallest key, ErrEmptySymbolTable if the symbol table is empty.
// The complexity is O(N) where N is the number of key-value pairs.
func (b *BST[K, V]) Min() (K, error) {
	if b.IsEmpty() {
		var key K
		return key, ErrEmptySymbolTable
	}
	return b.min(b.root).key, nil
}

func (b *BST[K, V]) min(node *bstNode[K, V]) *bstNode[K, V] {
	if node.left == nil {
		return node
	}
	return b.min(node.left)
}

// Max returns the largest key, ErrEmptySymbolTable if the symbol table is empty.
// The complexity is O(N) where N is the number of key-value pairs.
func (b *BST[K, V]) Max() (K, error) {
	if b.IsEmpty() {
		var key K
		return key, ErrEmptySymbolTable
	}
	return b.max(b.root).key, nil
}

func (b *BST[K, V]) max(node *bstNode[K, V]) *bstNode[K, V] {
	if node.right == nil {
		return node
	}
	return b.max(node.right)
}

// Delete removes the specified key and associated value, ErrAbsentKey if key is absent.
// The complexity is O(N) where N is the number of key-value pairs.
func (b *BST[K, V]) Delete(key K) error {
	root, err := b.delete(b.root, key)
	b.root = root
	return err
}

func (b *BST[K, V]) delete(node *bstNode[K, V], key K) (*bstNode[K, V], error) {
	if node == nil {
		return nil, ErrAbsentKey
	}

	var err error = nil

	cmp := b.compare(key, node.key)
	if cmp < 0 {
		node.left, err = b.delete(node.left, key)
	} else if cmp > 0 {
		node.right, err = b.delete(node.right, key)
	} else {
		if node.left == nil {
			return node.right, nil
		}
		if node.right == nil {
			return node.left, nil
		}
		tmpNode := node
		node = b.min(tmpNode.right)
		node.right = b.delMin(tmpNode.right)
		node.left = tmpNode.left
	}

	node.size = 1 + b.size(node.left) + b.size(node.right)
	return node, err
}

// Floor returns the largest key less than or equal to key, ErrTooSmallFloorKey if key to floor is too small.
// The complexity is O(N) where N is the number of key-value pairs.
func (b *BST[K, V]) Floor(key K) (K, error) {
	node, err := b.floor(b.root, key)
	if err != nil {
		return key, err
	}
	return node.key, nil
}

func (b *BST[K, V]) floor(node *bstNode[K, V], key K) (*bstNode[K, V], error) {
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
// The complexity is O(N) where N is the number of key-value pairs.
func (b *BST[K, V]) Ceiling(key K) (K, error) {
	node, err := b.ceiling(b.root, key)
	if err != nil {
		return key, err
	}
	return node.key, nil
}

func (b *BST[K, V]) ceiling(node *bstNode[K, V], key K) (*bstNode[K, V], error) {
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
// The complexity is O(N) where N is the number of key-value pairs.
func (b *BST[K, V]) Select(k int) (K, error) {
	var key K
	if k < 0 || k >= b.Size() {
		return key, ErrInvalidRank
	}
	return b.selectRecursive(b.root, k)
}

func (b *BST[K, V]) selectRecursive(node *bstNode[K, V], k int) (K, error) {
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
// The complexity is O(N) where N is the number of key-value pairs.
func (b *BST[K, V]) Rank(key K) (int, error) {
	return b.rank(b.root, key)
}

func (b *BST[K, V]) rank(node *bstNode[K, V], key K) (int, error) {
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
// It takes O(N) time to prepare iterator where N is the number of key-value pairs.
func (b *BST[K, V]) Iterator() iter.Seq2[K, V] {
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
// It takes O(N) time to prepare iterator where N is the number of key-value pairs.
func (b *BST[K, V]) RangeIterator(lo, hi K) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if b.compare(lo, hi) > 0 {
			return
		}
		b.iterator(yield, b.root, lo, hi)
	}
}

func (b *BST[K, V]) iterator(yield func(K, V) bool, node *bstNode[K, V], lo, hi K) {
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
// The complexity is O(N) where N is the number of key-value pairs.
func (b *BST[K, V]) RangeSize(lo, hi K) int {
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
