package graph

import (
	"errors"
	"github.com/inpour/algorithms/fundamental"
	"iter"
)

// Graph represents an undirected graph of vertices named 0 through v – 1. This implementation uses an adjacency-lists
// representation, which is a vertex-indexed array of Bags.
// Parallel edges and self-loops are permitted. By convention, a self-loop v-v appears in the adjacency list of v twice
// and contributes two to the degree of v.
// It uses O(V + E) space, where V is the number of vertices and E is the number of edges.
type Graph struct {
	v   int                     // number of vertices
	e   int                     // number of edges
	adj []*fundamental.Bag[int] // adjacent vertices
}

// NewGraph initializes a graph with v number vertices
// The complexity is O(V), where V is the number of vertices.
func NewGraph(v int) (*Graph, error) {
	if v < 0 {
		return nil, ErrInvalidVertices
	}

	adj := make([]*fundamental.Bag[int], v)
	for i := 0; i < v; i++ {
		adj[i] = fundamental.NewBag[int]()
	}

	return &Graph{
		v:   v,
		e:   0,
		adj: adj,
	}, nil
}

var ErrInvalidVertices = errors.New("number of vertices in a Graph must be non-negative")
var ErrInvalidVertexIndex = errors.New("invalid vertex index")

// V returns the number of vertices.
// The complexity is O(1).
func (graph *Graph) V() int {
	return graph.v
}

// E returns the number of edges.
// The complexity is O(1).
func (graph *Graph) E() int {
	return graph.e
}

func (graph *Graph) validateVertex(v int) error {
	if v < 0 || v >= graph.v {
		return ErrInvalidVertexIndex
	}
	return nil
}

// AddEdge adds the undirected edge v-w.
// The complexity is O(1).
func (graph *Graph) AddEdge(v, w int) error {
	if err := graph.validateVertex(v); err != nil {
		return err
	}
	if err := graph.validateVertex(w); err != nil {
		return err
	}
	graph.e++
	graph.adj[v].Add(w)
	graph.adj[w].Add(v)
	return nil
}

// Adj returns an iterator that iterates over vertices adjacent to vertex v.
func (graph *Graph) Adj(v int) (iter.Seq[int], error) {
	if err := graph.validateVertex(v); err != nil {
		return nil, err
	}
	return graph.adj[v].Iterator(), nil
}

// Degree returns the degree of vertex v.
// The complexity is O(1).
func (graph *Graph) Degree(v int) (int, error) {
	if err := graph.validateVertex(v); err != nil {
		return -1, err
	}
	return graph.adj[v].Size(), nil
}
