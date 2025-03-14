package graph

import (
	"github.com/inpour/algorithms/fundamental"
	"iter"
)

// BreadthFirstPath represents a data type for finding the shortest paths from a source vertex (s) to every other vertex in graph.
// This implementation uses breadth-first search (BFS).
// It uses O(V) extra space (not including the graph), where V is the number of vertices.
type BreadthFirstPath struct {
	marked []bool // marked[v] = is there an s-v path?
	edgeTo []int  // edgeTo[v] = last edge on shortest s-v path
	distTo []int  // distTo[v] = number of edges in shortest s-v path
}

// NewBreadthFirstPath computes the shortest path between the source vertex (s) and every other vertex in graph.
// The complexity is O(V + E), where V is the number of vertices and E is the number of edges.
func NewBreadthFirstPath(graph UndirectedOrDirectedGraph, s int) (*BreadthFirstPath, error) {
	if err := graph.validateVertex(s); err != nil {
		return nil, err
	}
	b := &BreadthFirstPath{
		marked: make([]bool, graph.V()),
		edgeTo: make([]int, graph.V()),
		distTo: make([]int, graph.V()),
	}
	b.bfs(graph, s)
	return b, nil
}

// NewBreadthFirstPathMultiSource computes the shortest path between any one of the source vertices and every other vertex in graph.
// The complexity is O(V + E), where V is the number of vertices and E is the number of edges.
func NewBreadthFirstPathMultiSource(graph UndirectedOrDirectedGraph, sources []int) (*BreadthFirstPath, error) {
	for i := 0; i < len(sources); i++ {
		if err := graph.validateVertex(sources[i]); err != nil {
			return nil, err
		}
	}
	b := &BreadthFirstPath{
		marked: make([]bool, graph.V()),
		edgeTo: make([]int, graph.V()),
		distTo: make([]int, graph.V()),
	}
	b.bfsMultiSource(graph, sources)
	return b, nil
}

// bfs (breadth first search) from a single source
func (b *BreadthFirstPath) bfs(graph UndirectedOrDirectedGraph, s int) {
	for v := 0; v < graph.V(); v++ {
		b.distTo[v] = -1
	}
	b.distTo[s] = 0
	b.marked[s] = true
	q := fundamental.NewQueue[int]()
	q.Enqueue(s)
	for !q.IsEmpty() {
		v, _ := q.Dequeue()
		adj, _ := graph.Adj(v)
		for w := range adj {
			if !b.marked[w] {
				b.edgeTo[w] = v
				b.distTo[w] = b.distTo[v] + 1
				b.marked[w] = true
				q.Enqueue(w)
			}
		}
	}
}

// bfsMultiSource (breadth first search) from multiple sources
func (b *BreadthFirstPath) bfsMultiSource(graph UndirectedOrDirectedGraph, sources []int) {
	for v := 0; v < graph.V(); v++ {
		b.distTo[v] = -1
	}
	q := fundamental.NewQueue[int]()
	for _, s := range sources {
		b.distTo[s] = 0
		b.marked[s] = true
		q.Enqueue(s)
	}
	for !q.IsEmpty() {
		v, _ := q.Dequeue()
		adj, _ := graph.Adj(v)
		for w := range adj {
			if !b.marked[w] {
				b.edgeTo[w] = v
				b.distTo[w] = b.distTo[v] + 1
				b.marked[w] = true
				q.Enqueue(w)
			}
		}
	}
}

// HasPathTo returns true if there is a path between the source vertex and vertex v.
// The complexity is O(1).
func (b *BreadthFirstPath) HasPathTo(v int) (bool, error) {
	if err := b.validateVertex(v); err != nil {
		return false, err
	}
	return b.marked[v], nil
}

// DistTo returns the number of edges in the shortest path between the source vertex and vertex v, or
// returns -1 if there is no path.
// The complexity is O(1).
func (b *BreadthFirstPath) DistTo(v int) (int, error) {
	if err := b.validateVertex(v); err != nil {
		return -1, err
	}
	return b.distTo[v], nil
}

// PathTo returns an iterator that iterates over the shortest path between the source vertex and vertex v.
func (b *BreadthFirstPath) PathTo(v int) (iter.Seq[int], error) {
	s := fundamental.NewStack[int]()
	if hasPath, err := b.HasPathTo(v); !hasPath {
		return s.Iterator(), err
	}
	x := v
	for ; b.distTo[x] != 0; x = b.edgeTo[x] {
		s.Push(x)
	}
	s.Push(x)
	return s.Iterator(), nil
}

func (b *BreadthFirstPath) validateVertex(v int) error {
	if v < 0 || v >= len(b.marked) {
		return ErrInvalidVertexIndex
	}
	return nil
}
