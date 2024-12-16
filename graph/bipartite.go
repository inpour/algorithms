package graph

import (
	"errors"
	"github.com/inpour/algorithms/fundamental"
	"iter"
)

// Bipartite represents a data type for determining whether an undirected graph is bipartite< or whether it has an odd-length cycle.
// A graph is bipartite if and only if it has no odd-length cycle.
// This implementation uses breadth-first search.
// It uses O(V) extra space (not including the graph), where V is the number of vertices.
type Bipartite struct {
	isBipartite bool
	color       []bool
	marked      []bool
	edgeTo      []int
	cycle       *fundamental.Queue[int]
}

// NewBipartite determines whether an undirected graph is bipartite and finds either a bipartition or an odd-length cycle.
// The complexity is O(V + E), where V is the number of vertices and E is the number of edges.
func NewBipartite(graph *Graph) *Bipartite {
	b := &Bipartite{
		isBipartite: true,
		color:       make([]bool, graph.V()),
		marked:      make([]bool, graph.V()),
		edgeTo:      make([]int, graph.V()),
		cycle:       fundamental.NewQueue[int](),
	}
	for v := 0; v < graph.V() && b.isBipartite; v++ {
		if !b.marked[v] {
			b.bfs(graph, v)
		}
	}
	return b
}

var ErrGraphIsNotBipartite = errors.New("graph is not bipartite")

// bfs (breadth first search) from s
func (b *Bipartite) bfs(graph *Graph, s int) {
	q := fundamental.NewQueue[int]()
	b.color[s] = false
	b.marked[s] = true
	q.Enqueue(s)

	for !q.IsEmpty() {
		v, _ := q.Dequeue()
		adj, _ := graph.Adj(v)
		for w := range adj {
			if !b.marked[w] {
				b.marked[w] = true
				b.edgeTo[w] = v
				b.color[w] = !b.color[v]
				q.Enqueue(w)
			} else if b.color[w] == b.color[v] {
				b.isBipartite = false

				// to form odd cycle, consider s-v path and s-w path
				// and let x be the closest node to v and w common to two paths
				// then (w-x path) + (x-v path) + (edge v-w) is an odd-length cycle
				// Note: distTo[v] == distTo[w];
				stack := fundamental.NewStack[int]()
				x := v
				y := w
				for x != y {
					stack.Push(x)
					b.cycle.Enqueue(y)
					x = b.edgeTo[x]
					y = b.edgeTo[y]
				}
				stack.Push(x)
				for !stack.IsEmpty() {
					x, _ = stack.Pop()
					b.cycle.Enqueue(x)
				}
				b.cycle.Enqueue(w)
				return
			}
		}
	}
}

// IsBipartite returns true if the graph is bipartite.
// The complexity is O(1).
func (b *Bipartite) IsBipartite() bool {
	return b.isBipartite
}

// Color returns the side of the bipartite that vertex v is on.
// The complexity is O(1).
func (b *Bipartite) Color(v int) (bool, error) {
	if err := b.validateVertex(v); err != nil {
		return false, err
	}
	if !b.isBipartite {
		return false, ErrGraphIsNotBipartite
	}

	return b.color[v], nil
}

// OddCycle returns an iterator that iterates over an odd-length cycle if the graph is not bipartite.
func (b *Bipartite) OddCycle() iter.Seq[int] {
	return b.cycle.Iterator()
}

func (b *Bipartite) validateVertex(v int) error {
	if v < 0 || v >= len(b.marked) {
		return ErrInvalidVertexIndex
	}
	return nil
}
