package graph

import (
	"github.com/inpour/algorithms/fundamental"
	"iter"
)

// DepthFirstPath represents a data type for finding paths from a source vertex (s) to every other vertex in graph.
// This implementation uses depth-first search (DFS).
type DepthFirstPath struct {
	marked []bool // marked[v] = is there an s-v path?
	edgeTo []int  // edgeTo[v] = last edge on s-v path
	s      int    // source vertex
}

// NewDepthFirstPath computes a path between s and every other vertex in graph.
// It uses Θ(v) extra space (not including the graph) and takes Θ(v+e) time in the worst case,
// where v is the number of vertices and e is the number of edges.
func NewDepthFirstPath(graph *Graph, s int) (*DepthFirstPath, error) {
	if err := graph.validateVertex(s); err != nil {
		return nil, err
	}
	d := &DepthFirstPath{
		marked: make([]bool, graph.V()),
		edgeTo: make([]int, graph.V()),
		s:      s,
	}
	d.dfs(graph, s)
	return d, nil
}

// dfs (depth first search) from v
func (d *DepthFirstPath) dfs(graph *Graph, v int) {
	d.marked[v] = true
	adj, _ := graph.Adj(v)
	for w := range adj {
		if !d.marked[w] {
			d.edgeTo[w] = v
			d.dfs(graph, w)
		}
	}
}

// HasPathTo returns true if there is a path between the source vertex and vertex v.
// The complexity is O(1).
func (d *DepthFirstPath) HasPathTo(v int) (bool, error) {
	if v < 0 || v >= len(d.marked) {
		return false, ErrInvalidVertexIndex
	}
	return d.marked[v], nil
}

// PathTo returns an iterator that iterates over a path between the source vertex and vertex v.
func (d *DepthFirstPath) PathTo(v int) (iter.Seq[int], error) {
	s := fundamental.NewStack[int]()
	if hasPath, err := d.HasPathTo(v); !hasPath {
		return s.Iterator(), err
	}
	for x := v; x != d.s; x = d.edgeTo[x] {
		s.Push(x)
	}
	s.Push(d.s)
	return s.Iterator(), nil
}
