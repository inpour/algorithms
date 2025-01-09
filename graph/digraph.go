package graph

import (
	"github.com/inpour/algorithms/fundamental"
	"iter"
)

// Digraph represents a directed graph of vertices named 0 through v – 1. This implementation uses an adjacency-lists
// representation, which is a vertex-indexed array of Bags.
// Parallel edges and self-loops are permitted.
// It uses O(V + E) space, where V is the number of vertices and E is the number of edges.
type Digraph struct {
	v        int                     // number of vertices
	e        int                     // number of edges
	adj      []*fundamental.Bag[int] // adjacent vertices
	inDegree []int                   // inDegree[v] = in-degree of vertex v
}

// NewDigraph initializes a graph with v number vertices
// The complexity is O(V), where V is the number of vertices.
func NewDigraph(v int) (*Digraph, error) {
	if v < 0 {
		return nil, ErrInvalidVertices
	}

	adj := make([]*fundamental.Bag[int], v)
	for i := 0; i < v; i++ {
		adj[i] = fundamental.NewBag[int]()
	}

	return &Digraph{
		v:        v,
		e:        0,
		adj:      adj,
		inDegree: make([]int, v),
	}, nil
}

// V returns the number of vertices.
// The complexity is O(1).
func (digraph *Digraph) V() int {
	return digraph.v
}

// E returns the number of edges.
// The complexity is O(1).
func (digraph *Digraph) E() int {
	return digraph.e
}

func (digraph *Digraph) validateVertex(v int) error {
	if v < 0 || v >= digraph.v {
		return ErrInvalidVertexIndex
	}
	return nil
}

// AddEdge adds the directed edge v-w.
// The complexity is O(1).
func (digraph *Digraph) AddEdge(v, w int) error {
	if err := digraph.validateVertex(v); err != nil {
		return err
	}
	if err := digraph.validateVertex(w); err != nil {
		return err
	}
	digraph.e++
	digraph.adj[v].Add(w)
	digraph.inDegree[w]++
	return nil
}

// Adj returns an iterator that iterates over vertices adjacent to vertex v.
// The complexity is O(1) (Though, iterating over the vertices returned by Adj(v) takes time proportional to the
// out-degree of the vertex v).
func (digraph *Digraph) Adj(v int) (iter.Seq[int], error) {
	if err := digraph.validateVertex(v); err != nil {
		return nil, err
	}
	return digraph.adj[v].Iterator(), nil
}

// InDegree returns the in-degree of vertex v.
// The complexity is O(1).
func (digraph *Digraph) InDegree(v int) (int, error) {
	if err := digraph.validateVertex(v); err != nil {
		return -1, err
	}
	return digraph.inDegree[v], nil
}

// OutDegree returns the out-degree of vertex v.
// The complexity is O(1).
func (digraph *Digraph) OutDegree(v int) (int, error) {
	if err := digraph.validateVertex(v); err != nil {
		return -1, err
	}
	return digraph.adj[v].Size(), nil
}

// Reverse returns the reverse of the digraph.
// The complexity is O(V + E), where V is the number of vertices and E is the number of edges.
func (digraph *Digraph) Reverse() *Digraph {
	reverse, _ := NewDigraph(digraph.v)
	for v := 0; v < digraph.v; v++ {
		adj, _ := digraph.Adj(v)
		for w := range adj {
			reverse.AddEdge(w, v)
		}
	}
	return reverse
}
