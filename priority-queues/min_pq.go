package priority_queues

// MinPQ represents min priority queue of generic key.
// It relies on the compare() function to compare two keys:
//
//	if a == b then compare(a, b) returns 0
//	if a > b then compare(a, b) returns 1
//	if a < b then compare(a, b) returns -1
//
// This implementation uses a min heap as underlying data structure.
type MinPQ[K any] struct {
	pq      []K              // store items at indices 1 to n
	n       int              // number of items on priority queue
	compare func(a, b K) int // function to compare two keys
}

// NewMinPQ initializes an empty min priority queue.
// It gets a function as a parameter to compare two keys.
// The complexity is O(1).
func NewMinPQ[K any](compare func(a, b K) int) *MinPQ[K] {
	return &MinPQ[K]{
		pq:      make([]K, 2),
		n:       0,
		compare: compare,
	}
}

// IsEmpty returns true if priority queue is empty.
// The complexity is O(1).
func (m *MinPQ[K]) IsEmpty() bool {
	return m.n == 0
}

// Size returns the size priority queue is empty.
// The complexity is O(1).
func (m *MinPQ[K]) Size() int {
	return m.n
}

// Min returns the smallest key on this priority queue.
// The complexity is O(1).
func (m *MinPQ[K]) Min() (K, error) {
	var key K
	if m.IsEmpty() {
		return key, ErrEmptyPQ
	}
	return m.pq[1], nil
}

// Insert adds a new key to this priority queue.
// The complexity is O(log(N)) where N is the number of keys in priority queue.
func (m *MinPQ[K]) Insert(key K) {
	if m.n == len(m.pq)-1 {
		m.resize(2 * len(m.pq))
	}
	m.n++
	m.pq[m.n] = key
	m.swim(m.n)
}

// DelMin removes and returns the smallest key on this priority queue.
// The complexity is O(log(N)) where N is the number of keys in priority queue.
func (m *MinPQ[K]) DelMin() (K, error) {
	var key K
	if m.IsEmpty() {
		return key, ErrEmptyPQ
	}
	min_ := m.pq[1]
	m.exchange(1, m.n)
	m.n--
	m.sink(1)
	m.pq[m.n+1] = key // to avoid loitering and help garbage collection
	if (m.n > 0) && (m.n == (len(m.pq)-1)/4) {
		m.resize(len(m.pq) / 2)
	}
	return min_, nil
}

// resize the underlying slice
func (m *MinPQ[K]) resize(newSize int) {
	if newSize <= m.n {
		return
	}
	temp := make([]K, newSize)
	for i := 1; i <= m.n; i++ {
		temp[i] = m.pq[i]
	}
	m.pq = temp
}

func (m *MinPQ[K]) swim(i int) {
	for i > 1 && m.greater(i/2, i) {
		m.exchange(i/2, i)
		i = i / 2
	}
}

func (m *MinPQ[K]) sink(i int) {
	for 2*i <= m.n {
		j := 2 * i
		if j < m.n && m.greater(j, j+1) {
			j++
		}
		if !m.greater(i, j) {
			break
		}
		m.exchange(i, j)
		i = j
	}
}

func (m *MinPQ[K]) exchange(i, j int) {
	swap := m.pq[i]
	m.pq[i] = m.pq[j]
	m.pq[j] = swap
}

func (m *MinPQ[K]) greater(i, j int) bool {
	return m.compare(m.pq[i], m.pq[j]) > 0
}
