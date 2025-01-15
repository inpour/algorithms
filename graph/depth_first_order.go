package graph

import (
	"github.com/inpour/algorithms/fundamental"
	"iter"
)

// DepthFirstOrder represents a data type for determining depth-first search ordering of the vertices in a digraph.
// This implementation uses depth-first search.
// It uses O(V) extra space (not including the digraph), where V is the number of vertices.
type DepthFirstOrder struct {
	marked      []bool                  // marked[v] = has v been marked in dfs?
	pre         *fundamental.Queue[int] // vertices in preorder
	post        *fundamental.Queue[int] // vertices in postorder
	reversePost *fundamental.Stack[int] // vertices in reverse postorder
}

// NewDepthFirstOrder determines a depth-first order for the digraph.
// The complexity is O(V + E), where V is the number of vertices and E is the number of edges.
func NewDepthFirstOrder(digraph *Digraph) *DepthFirstOrder {
	d := &DepthFirstOrder{
		marked:      make([]bool, digraph.V()),
		pre:         fundamental.NewQueue[int](),
		post:        fundamental.NewQueue[int](),
		reversePost: fundamental.NewStack[int](),
	}
	for v := 0; v < digraph.V(); v++ {
		if !d.marked[v] {
			d.dfs(digraph, v)
		}
	}
	return d
}

// dfs (depth first search) from v
func (d *DepthFirstOrder) dfs(digraph *Digraph, v int) {
	d.pre.Enqueue(v)
	d.marked[v] = true
	adj, _ := digraph.Adj(v)
	for w := range adj {
		if !d.marked[w] {
			d.dfs(digraph, w)
		}
	}
	d.post.Enqueue(v)
	d.reversePost.Push(v)
}

// Pre returns the vertices in preorder, as an iterable of vertices.
// The complexity is O(1).
func (d *DepthFirstOrder) Pre() iter.Seq[int] {
	return d.pre.Iterator()
}

// Post returns the vertices in postorder, as an iterable of vertices.
// The complexity is O(1).
func (d *DepthFirstOrder) Post() iter.Seq[int] {
	return d.post.Iterator()
}

// ReversePost returns the vertices in reverse postorder, as an iterable of vertices.
// The complexity is O(1).
func (d *DepthFirstOrder) ReversePost() iter.Seq[int] {
	return d.reversePost.Iterator()
}

func (d *DepthFirstOrder) validateVertex(v int) error {
	if v < 0 || v >= len(d.marked) {
		return ErrInvalidVertexIndex
	}
	return nil
}
