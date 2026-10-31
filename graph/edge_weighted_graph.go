package graph

import (
	"iter"

	"github.com/inpour/algorithms/fundamental"
)

// EdgeWeightedGraph class represents an edge-weighted graph of vertices named 0 through V – 1, where each
// undirected edge is of type WeightedEdge and has a real-valued weight. It supports the following two primary
// operations: add an edge to the graph, iterate over all the edges incident with a vertex. It also provides
// methods for returning the degree of a vertex, the number of vertices V in the graph, and the number of edges E
// in the graph. Parallel edges and self-loops are permitted. By convention, a self-loop v-v appears in the
// adjacency list of v twice and contributes two to the degree of v.
// This implementation uses an adjacency-lists representation, which is a vertex-indexed array of Bag objects.
// It uses Θ(V + E) space, where E is the number of edges and V is the number of vertices.
// All instance methods take Θ(1) time.
// (Though, iterating over the edges takes time proportional to the degree of the vertex.)
type EdgeWeightedGraph struct {
	v   int                              // number of vertices
	e   int                              // number of edges
	adj []*fundamental.Bag[WeightedEdge] // adjacent vertices
}

// NewEdgeWeightedGraph initializes an empty edge-weighted graph with {@code V} vertices and 0 edges.
func NewEdgeWeightedGraph(v int) (*EdgeWeightedGraph, error) {
	if v < 0 {
		return nil, ErrInvalidVertices
	}

	adj := make([]*fundamental.Bag[WeightedEdge], v)
	for i := range v {
		adj[i] = fundamental.NewBag[WeightedEdge]()
	}

	return &EdgeWeightedGraph{
		v:   v,
		e:   0,
		adj: adj,
	}, nil
}

// V returns the number of vertices in this edge-weighted graph.
func (g *EdgeWeightedGraph) V() int {
	return g.v
}

// E returns the number of edges in this edge-weighted graph.
func (g *EdgeWeightedGraph) E() int {
	return g.e
}

func (g *EdgeWeightedGraph) validateVertex(v int) error {
	if v < 0 || v >= g.v {
		return ErrInvalidVertexIndex
	}
	return nil
}

// AddEdge adds the undirected edge e to this edge-weighted graph.
// The complexity is O(1).
func (g *EdgeWeightedGraph) AddEdge(e WeightedEdge) error {
	v := e.Either()
	w, _ := e.Other(v)
	if err := g.validateVertex(v); err != nil {
		return err
	}
	if err := g.validateVertex(w); err != nil {
		return err
	}
	g.e++
	g.adj[v].Add(e)
	g.adj[w].Add(e)
	return nil
}

// Adj returns an iterator over edges incident to vertex v.
// The complexity is O(1) (Though, iterating over the vertices returned by Adj(v) takes time proportional to the
// degree of the vertex v).
func (g *EdgeWeightedGraph) Adj(v int) (iter.Seq[WeightedEdge], error) {
	if err := g.validateVertex(v); err != nil {
		return nil, err
	}
	return g.adj[v].Iterator(), nil
}

// Degree returns the degree of vertex v.
// The complexity is O(1).
func (g *EdgeWeightedGraph) Degree(v int) (int, error) {
	if err := g.validateVertex(v); err != nil {
		return -1, err
	}
	return g.adj[v].Size(), nil
}
