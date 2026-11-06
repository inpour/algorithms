package priority_queues

import (
	"errors"
)

// IndexMaxPQ represents an indexed max priority queue of generic key.
//
// While in a standard priority queue only insert items and remove the maximum is possible, an
// indexed priority queue can associate each key with an integer index (from 0 to maxN - 1).
// This allows to directly reference, update, decrease, increase, or delete specific keys efficiently
// without needing to search through the entire queue.
// In a standard priority queue, if an element's priority changes, have no way to find that element
// inside the heap in less than O(n) time. IndexMaxPQ solves this by maintaining a direct mapping between
// external indices and their positions in the heap, bringing update operations down to O(log n) time.
//
// It relies on the compare() function to compare two keys:
//
//	if a == b then compare(a, b) returns 0
//	if a > b then compare(a, b) returns 1
//	if a < b then compare(a, b) returns -1
//
// In order to let the client refer to keys on the priority queue, an integer between 0 and maxN - 1
// is associated with each key—the client uses this integer to specify which key to delete or change.
//
// This implementation uses a binary heap along with an array to associate keys with integers in the given range.
type IndexMaxPQ[K any] struct {
	maxN    int              // maximum number of elements on PQ
	n       int              // number of items on priority queue
	pq      []int            // binary heap using 1-based indexing
	qp      []int            // inverse of pq (qp[pq[i]] = pq[qp[i]] = i)
	keys    []K              // keys[i] = priority of i (stores the actual priority keys)
	compare func(a, b K) int // function to compare two keys
}

// NewIndexMaxPQ initializes an indexed priority queue with indices between 0 and maxN - 1.
// It gets a function as a parameter to compare two keys.
// The complexity is O(1).
func NewIndexMaxPQ[K any](compare func(a, b K) int, maxN int) (*IndexMaxPQ[K], error) {
	if maxN < 0 {
		return nil, errors.New("maxN must be non-negative")
	}

	qp := make([]int, maxN+1)
	for i := range maxN + 1 {
		qp[i] = -1
	}

	return &IndexMaxPQ[K]{
		maxN:    maxN,
		n:       0,
		pq:      make([]int, maxN+1),
		qp:      qp,
		keys:    make([]K, maxN+1),
		compare: compare,
	}, nil
}

// IsEmpty returns true if priority queue is empty.
// The complexity is O(1).
func (p *IndexMaxPQ[K]) IsEmpty() bool {
	return p.n == 0
}

func (p *IndexMaxPQ[K]) validateIndex(i int) error {
	if i < 0 || i >= p.maxN {
		return ErrInexOutOfRange
	}
	return nil
}

// Contains returns true if i is an index of this PQ.
// The complexity is O(1).
func (p *IndexMaxPQ[K]) Contains(i int) (bool, error) {
	if err := p.validateIndex(i); err != nil {
		return false, err
	}
	return p.qp[i] != -1, nil
}

// Size returns the size priority queue is empty.
// The complexity is O(1).
func (p *IndexMaxPQ[K]) Size() int {
	return p.n
}

// Insert adds a new index and associated key to this priority queue.
// The complexity is O(log(N)) where N is the number of keys in priority queue.
func (p *IndexMaxPQ[K]) Insert(i int, key K) error {
	contains, err := p.Contains(i)
	if err != nil {
		return err
	}
	if contains {
		return ErrInexAlreadyExists
	}
	p.n++
	p.qp[i] = p.n
	p.pq[p.n] = i
	p.keys[i] = key
	p.swim(p.n)
	return nil
}

// MaxIndex returns an index associated with a maximum key.
// The complexity is O(1).
func (p *IndexMaxPQ[K]) MaxIndex() (int, error) {
	if p.IsEmpty() {
		return -1, ErrEmptyPQ
	}
	return p.pq[1], nil
}

// MaxKey returns a maximum key.
// The complexity is O(1).
func (p *IndexMaxPQ[K]) MaxKey() (K, error) {
	if p.IsEmpty() {
		var k K
		return k, ErrEmptyPQ
	}
	return p.keys[p.pq[1]], nil
}

// DelMax removes a maximum key and returns its associated index.
// The complexity is O(log(N)) where N is the number of keys in priority queue.
func (p *IndexMaxPQ[K]) DelMax() (int, error) {
	if p.IsEmpty() {
		return -1, ErrEmptyPQ
	}
	max_ := p.pq[1]
	p.exchange(1, p.n)
	p.n--
	p.sink(1)
	p.qp[max_] = -1 // delete
	var key K
	p.keys[max_] = key // to help with garbage collection
	p.pq[p.n+1] = -1   // not needed
	return max_, nil
}

// KeyOf returns the key associated with index i.
// The complexity is O(1).
func (p *IndexMaxPQ[K]) KeyOf(i int) (K, error) {
	var key K
	contains, err := p.Contains(i)
	if err != nil {
		return key, err
	}
	if !contains {
		return key, ErrInexNotExists
	}
	return p.keys[i], nil
}

// ChangeKey changes the key associated with index {@code i} to the specified value.
// The complexity is O(log(N)) where N is the number of keys in priority queue.
func (p *IndexMaxPQ[K]) ChangeKey(i int, key K) error {
	contains, err := p.Contains(i)
	if err != nil {
		return err
	}
	if !contains {
		return ErrInexNotExists
	}
	p.keys[i] = key
	p.swim(p.qp[i])
	p.sink(p.qp[i])
	return nil
}

// DecreaseKey decrease the key associated with index {@code i} to the specified value.
// The complexity is O(log(N)) where N is the number of keys in priority queue.
func (p *IndexMaxPQ[K]) DecreaseKey(i int, key K) error {
	contains, err := p.Contains(i)
	if err != nil {
		return err
	}
	if !contains {
		return ErrInexNotExists
	}
	if p.compare(p.keys[i], key) <= 0 {
		return ErrDecreaseWithInvalidKey
	}
	p.keys[i] = key
	p.sink(p.qp[i]) // In a MaxPQ, decreasing a key means it might need to sink
	return nil
}

// IncreaseKey increase the key associated with index {@code i} to the specified value.
// The complexity is O(log(N)) where N is the number of keys in priority queue.
func (p *IndexMaxPQ[K]) IncreaseKey(i int, key K) error {
	contains, err := p.Contains(i)
	if err != nil {
		return err
	}
	if !contains {
		return ErrInexNotExists
	}
	if p.compare(p.keys[i], key) >= 0 {
		return ErrIncreaseWithInvalidKey
	}
	p.keys[i] = key
	p.swim(p.qp[i]) // In a MaxPQ, increasing a key means it might need to swim
	return nil
}

// Delete remove the key associated with index i.
// The complexity is O(log(N)) where N is the number of keys in priority queue.
func (p *IndexMaxPQ[K]) Delete(i int) error {
	contains, err := p.Contains(i)
	if err != nil {
		return err
	}
	if !contains {
		return ErrInexNotExists
	}
	index := p.qp[i]
	p.exchange(index, p.n)
	p.n--
	p.swim(index)
	p.sink(index)
	var key K
	p.keys[i] = key
	p.qp[i] = -1
	return nil
}

func (p *IndexMaxPQ[K]) swim(i int) {
	for i > 1 && p.less(i/2, i) {
		p.exchange(i/2, i)
		i = i / 2
	}
}

func (p *IndexMaxPQ[K]) sink(i int) {
	for 2*i <= p.n {
		j := 2 * i
		if j < p.n && p.less(j, j+1) {
			j++
		}
		if !p.less(i, j) {
			break
		}
		p.exchange(i, j)
		i = j
	}
}

func (p *IndexMaxPQ[K]) exchange(i, j int) {
	swap := p.pq[i]
	p.pq[i] = p.pq[j]
	p.pq[j] = swap
	p.qp[p.pq[i]] = i
	p.qp[p.pq[j]] = j
}

func (p *IndexMaxPQ[K]) less(i, j int) bool {
	return p.compare(p.keys[p.pq[i]], p.keys[p.pq[j]]) < 0
}
