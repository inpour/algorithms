package fundamental

import (
	"errors"
	"sync"
)

// The UnionFind type represents a union–find data type (also known as the disjoint-sets
// data type). This implementation uses weighted quick union by subtreeSize (without path compression).
//
// It supports the classic union and find operations, along with a count operation that
// returns the total number of sets.
//
// The union–find data type models a collection of sets containing n elements, with each element in
// exactly one set. The elements are named 0 through n–1.
// Initially, there are n sets, with each element in its own set. The canonical element of a set (also
// known as the root, identifier, leader, or set representative) is one distinguished element in the set.
//
// Here is a summary of the operations:
//
//	find(p) returns the canonical element of the set containing p. The find operation returns the same
//	value for two elements if and only if they are in the same set.
//	union(p, q) merges the set containing element p with the set containing element q. That is, if
//	p and q are in different sets, replace these two sets with a new set that is the union of the two.
type UnionFind struct {
	lock        *sync.Mutex // protect race condition
	parent      []int       // parent[i] = parent of i (if parent[i] = i then i is root)
	subtreeSize []int       // subtreeSize[i] = number of elements in subtree rooted at i
	count       int         // number of sets
}

// NewUnionFind initializes an empty union-find data structure with n elements 0 through n-1.
// Initially, each element is in its own set.
// The complexity is O(N) where N = n = uf.Size().
func NewUnionFind(n int) *UnionFind {
	parent := make([]int, n)
	subtreeSize := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		subtreeSize[i] = 1
	}
	return &UnionFind{
		lock:        &sync.Mutex{},
		parent:      parent,
		subtreeSize: subtreeSize,
		count:       n,
	}
}

var ErrInvalidIndex = errors.New("invalid index")

// Count returns the number of sets.
// The complexity is O(1).
func (uf *UnionFind) Count() int {
	return uf.count
}

// Size returns the number of elements.
// The complexity is O(1).
func (uf *UnionFind) Size() int {
	return len(uf.parent)
}

// Find returns the canonical element of the set containing element p.
// The complexity is O(log(N)) where N = uf.Size().
func (uf *UnionFind) Find(p int) (int, error) {
	if err := uf.validate(p); err != nil {
		return -1, err
	}

	// traverse until find the root
	for p != uf.parent[p] {
		p = uf.parent[p]
	}
	return p, nil
}

// Connected returns true if the two elements are in the same set.
// The complexity is O(log(N)) where N = uf.Size().
func (uf *UnionFind) Connected(p, q int) bool {
	rootP, _ := uf.Find(p)
	rootQ, _ := uf.Find(q)
	return rootP == rootQ
}

// Union Merges the set containing element p with the set containing element q.
// The complexity is O(log(N)) where N = uf.Size().
func (uf *UnionFind) Union(p, q int) {
	uf.lock.Lock()
	defer uf.lock.Unlock()

	rootP, errP := uf.Find(p)
	rootQ, errQ := uf.Find(q)
	if errP != nil || errQ != nil || rootP == rootQ {
		return
	}

	// make smaller root point to larger one
	if uf.subtreeSize[rootP] < uf.subtreeSize[rootQ] {
		uf.parent[rootP] = rootQ
		uf.subtreeSize[rootQ] += uf.subtreeSize[rootP]
	} else {
		uf.parent[rootQ] = rootP
		uf.subtreeSize[rootP] += uf.subtreeSize[rootQ]
	}

	uf.count--
}

// validate that p is a valid index
func (uf *UnionFind) validate(p int) error {
	if p < 0 || p >= uf.Size() {
		return ErrInvalidIndex
	}
	return nil
}
