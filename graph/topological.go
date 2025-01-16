package graph

import (
	"errors"
	"iter"
)

// Topological represents a data type for determining a topological order of a directed acyclic graph (DAG).
// A digraph has a topological order if and only if it is a DAG.
// This implementation uses depth-first search.
// It uses O(V) extra space (not including the digraph), where V is the number of vertices.
type Topological struct {
	order iter.Seq[int] // iterable of vertices in topological order
	rank  []int         // rank[v] = rank of vertex v in order
}

// NewTopological determines whether the digraph has a topological order and, if so, finds such a topological order.
// The complexity is O(V + E), where V is the number of vertices and E is the number of edges.
func NewTopological(digraph *Digraph) *Topological {
	t := &Topological{
		order: nil,
		rank:  make([]int, digraph.V()),
	}
	finder := NewDirectedCycle(digraph)
	if !finder.HasCycle() {
		dfs := NewDepthFirstOrder(digraph)
		t.order = dfs.ReversePost()
		i := 0
		for v := range t.order {
			t.rank[v] = i
			i++
		}
	}
	return t
}

var ErrNotDAG = errors.New("digraph is not a DAG")

// HasOrder returns true if the digraph has a topological order (or equivalently, if the digraph is a DAG).
// The complexity is O(1).
func (t *Topological) HasOrder() bool {
	return t.order != nil
}

// Rank returns the rank of vertex v in a topological order of the digraph, ErrNotDAG if the digraph is not a DAG.
// The complexity is O(1).
func (t *Topological) Rank(v int) (int, error) {
	if !t.HasOrder() {
		return -1, ErrNotDAG
	}
	if err := t.validateVertex(v); err != nil {
		return -1, err
	}
	return t.rank[v], nil
}

// Order returns a topological order of the vertices (as an iterable), ErrNotDAG if the digraph is not a DAG.
// The complexity is O(1).
func (t *Topological) Order() (iter.Seq[int], error) {
	if !t.HasOrder() {
		return nil, ErrNotDAG
	}
	return t.order, nil
}

func (t *Topological) validateVertex(v int) error {
	if v < 0 || v >= len(t.rank) {
		return ErrInvalidVertexIndex
	}
	return nil
}
