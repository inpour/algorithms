package graph

import (
	"github.com/inpour/algorithms/fundamental"
	"iter"
)

// DirectedCycle represents a data type for determining whether a directed graph has a directed cycle.
// This implementation uses depth-first search (DFS).
// It uses O(V) extra space (not including the graph), where V is the number of vertices.
type DirectedCycle struct {
	marked  []bool                  // marked[v] = has vertex v been marked?
	edgeTo  []int                   // edgeTo[v] = previous vertex on path to v
	onStack []bool                  // onStack[v] = is vertex on the stack?
	cycle   *fundamental.Stack[int] // directed cycle
}

// NewDirectedCycle determines whether the digraph has a directed cycle and, if so, finds such a cycle.
// The complexity is O(V + E), where V is the number of vertices and E is the number of edges.
func NewDirectedCycle(digraph *Digraph) *DirectedCycle {
	d := &DirectedCycle{
		marked:  make([]bool, digraph.V()),
		edgeTo:  make([]int, digraph.V()),
		onStack: make([]bool, digraph.V()),
		cycle:   fundamental.NewStack[int](),
	}

	for v := 0; v < digraph.V(); v++ {
		if !d.marked[v] && d.cycle.IsEmpty() {
			d.dfs(digraph, v)
		}
	}
	return d
}

// dfs (depth first search)
func (d *DirectedCycle) dfs(digraph *Digraph, v int) {
	d.onStack[v] = true
	d.marked[v] = true
	adj, _ := digraph.Adj(v)
	for w := range adj {

		// short circuit if cycle already found
		if !d.cycle.IsEmpty() {
			return
		}

		// found a new vertex, then recur, otherwise trace back directed cycle
		if !d.marked[w] {
			d.edgeTo[w] = v
			d.dfs(digraph, w)
		} else if d.onStack[w] {
			for x := v; x != w; x = d.edgeTo[x] {
				d.cycle.Push(x)
			}
			d.cycle.Push(w)
			d.cycle.Push(v)
		}
	}
	d.onStack[v] = false
}

// HasCycle returns true if the digraph has a directed cycle.
// The complexity is O(1).
func (d *DirectedCycle) HasCycle() bool {
	return !d.cycle.IsEmpty()
}

// Cycle returns an iterator that iterates over a directed cycle in the digraph.
// The complexity is O(1).
func (d *DirectedCycle) Cycle() iter.Seq[int] {
	return d.cycle.Iterator()
}
