package priority_queues

import (
	"errors"
)

var ErrEmptyPQ = errors.New("priority queue is empty")

// MaxPQ represents max priority queue of generic key.
// It relies on the compare() function to compare two keys:
//
//	if a == b then compare(a, b) returns 0
//	if a > b then compare(a, b) returns 1
//	if a < b then compare(a, b) returns -1
//
// This implementation uses a max heap as underlying data structure.
type MaxPQ[K any] struct {
	pq      []K              // store items at indices 1 to n
	n       int              // number of items on priority queue
	compare func(a, b K) int // function to compare two keys
}

// NewMaxPQ initializes an empty max priority queue.
// It gets a function as a parameter to compare two keys.
// The complexity is O(1).
func NewMaxPQ[K any](compare func(a, b K) int) *MaxPQ[K] {
	return &MaxPQ[K]{
		pq:      make([]K, 2),
		n:       0,
		compare: compare,
	}
}

// IsEmpty returns true if priority queue is empty.
// The complexity is O(1).
func (m *MaxPQ[K]) IsEmpty() bool {
	return m.n == 0
}

// Size returns true if priority queue is empty.
// The complexity is O(1).
func (m *MaxPQ[K]) Size() int {
	return m.n
}

// Max returns the largest key on this priority queue.
// The complexity is O(1).
func (m *MaxPQ[K]) Max() (K, error) {
	var key K
	if m.IsEmpty() {
		return key, ErrEmptyPQ
	}
	return m.pq[1], nil
}

// Insert adds a new key to this priority queue.
// The complexity is O(log(N)) where N is the number of keys in priority queue.
func (m *MaxPQ[K]) Insert(key K) {
	if m.n == len(m.pq)-1 {
		m.resize(2 * len(m.pq))
	}
	m.n++
	m.pq[m.n] = key
	m.swim(m.n)
}

// DelMax removes and returns the largest key on this priority queue.
// The complexity is O(log(N)) where N is the number of keys in priority queue.
func (m *MaxPQ[K]) DelMax() (K, error) {
	var key K
	if m.IsEmpty() {
		return key, ErrEmptyPQ
	}
	max_ := m.pq[1]
	m.exchange(1, m.n)
	m.n--
	m.sink(1)
	m.pq[m.n+1] = key // to avoid loitering and help garbage collection
	if (m.n > 0) && (m.n == (len(m.pq)-1)/4) {
		m.resize(len(m.pq) / 2)
	}
	return max_, nil
}

// resize the underlying slice
func (m *MaxPQ[K]) resize(newSize int) {
	if newSize <= m.n {
		return
	}
	temp := make([]K, newSize)
	for i := 1; i <= m.n; i++ {
		temp[i] = m.pq[i]
	}
	m.pq = temp
}

func (m *MaxPQ[K]) swim(k int) {
	for k > 1 && m.less(k/2, k) {
		m.exchange(k/2, k)
		k = k / 2
	}
}

func (m *MaxPQ[K]) sink(k int) {
	for 2*k <= m.n {
		j := 2 * k
		if j < m.n && m.less(j, j+1) {
			j++
		}
		if !m.less(k, j) {
			break
		}
		m.exchange(k, j)
		k = j
	}
}

func (m *MaxPQ[K]) exchange(i, j int) {
	swap := m.pq[i]
	m.pq[i] = m.pq[j]
	m.pq[j] = swap
}

func (m *MaxPQ[K]) less(i, j int) bool {
	return m.compare(m.pq[i], m.pq[j]) < 0
}
